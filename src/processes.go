package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

type taskProcess struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Err    error  `json:"err"`
	URL    string `json:"url"`
}

type allProcesses struct {
	List []*taskProcess
}

const (
	procStatRunning  = "running"
	procStatError    = "error"
	procStatStopped  = "stopped"
	procStatBuilding = "building"
	procStatStopping = "stopping"
)

var allProcs = &allProcesses{}

func getProc(id string) *taskProcess {
	for _, proc := range allProcs.List {
		if proc.ID == id {
			return proc
		}
	}
	return nil
}

func createProc(payload postData) *taskProcess {
	var id string

	if proc := getProc(payload.ID); proc == nil {
		id = proc.ID
		proc.Stop()
	} else {
		id = payload.ID
	}

	proc := &taskProcess{
		ID:   id,
		Name: payload.Name,
		Type: payload.Type,
		URL:  payload.GitURL + payload.Compose}

	allProcs.List = append(allProcs.List, proc)

	return proc
}

func (proc *taskProcess) Run() {

	var cmd *exec.Cmd

	proc.Status = procStatBuilding

	switch proc.Type {
	case "docker":
		cmd = proc.DockerRun()
	}

	proc.Status = procStatRunning
	err := cmd.Run()
	checkAPI(err, proc)
	proc.Status = procStatStopped
}

func (proc *taskProcess) Stop() {

	var cmd *exec.Cmd

	switch proc.Type {
	case "docker":
		cmd = proc.DockerStop()
	}

	proc.Status = procStatStopping
	err := cmd.Run()
	checkAPI(err, proc)
	proc.Status = procStatStopped
}

func (proc *taskProcess) DockerStop() *exec.Cmd {
	dockerloc := "./docker/" + proc.ID
	composeloc := dockerloc + "/docker-compose.yml"

	// Stopping the container
	cmd := exec.Command("docker-compose", "-f", composeloc, "down")

	return cmd
}

func (proc *taskProcess) DockerRun() *exec.Cmd {
	id := proc.ID
	url := proc.URL
	dockerloc := "./docker/" + id
	logsloc := "./logs/" + id
	composeloc := dockerloc + "/docker-compose.yml"

	// Make url, read the compose file
	resp, err := http.Get(url)
	checkAPI(err, proc)
	defer resp.Body.Close()
	file, err := ioutil.ReadAll(resp.Body)
	checkAPI(err, proc)

	// Make directory for docker and logs and save file
	os.Mkdir(dockerloc, 0777)
	os.Mkdir(logsloc, 0777)
	err = ioutil.WriteFile(composeloc, file, 0777)
	checkAPI(err, proc)

	proc.Status = procStatBuilding
	// Build Commands
	err = exec.Command("docker-compose", "-f", composeloc, "down").Run()
	checkAPI(err, proc)
	err = exec.Command("docker-compose", "-f", composeloc, "pull").Run()
	checkAPI(err, proc)

	// Get logging Running
	cmd := exec.Command("docker-compose", "-f", composeloc, "up", "--no-color")

	outfile, err := os.Create(logsloc + "/info.log")
	checkAPI(err, proc)
	defer outfile.Close()
	cmd.Stdout = outfile

	return cmd
}
