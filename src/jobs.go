package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

// Pointer to variable storing all jobs
var allJobs = &allTaskJobs{}

func getJob(id string) *taskJob {
	// If there are no jobs to chose from
	if len(allJobs.List) == 0 {
		return nil
	}

	// Return a job with the matching ID
	for _, job := range allJobs.List {
		if job.ID == id {
			return job
		}
	}
	return nil
}

// Creates and Runs the job based off of the payload
func createJob(payload postData) {

	// Stores info on the job from the task
	newJob := taskJob{
		Name: payload.Name,
		Type: payload.Type,
		URL:  payload.GitURL + payload.Compose,
		Enviroment: payload.Enviroment}

	// If a job with that ID already exists
	if job := getJob(payload.ID); job != nil {

		// Assign its ID to the old job ID
		newJob.ID = job.ID

		// Stop the old job
		job.Stop()

		// Replace and run the new job
		*job = newJob
		job.Run()

	} else {
		newJob.ID = payload.ID
		allJobs.List = append(allJobs.List, &newJob)
		comm := &newJob
		comm.Run()
	}
}

// Run a Job
func (job *taskJob) Run() {

	var cmd *exec.Cmd

	// Mark as building
	job.Status = jobStatus.Building

	// Run the correct job Type
	switch job.Type {
	case "docker":
		cmd = job.DockerCmd()
	}

	// If setting up the command went fine
	if job.Status != jobStatus.Errored {

		// Run Command
		job.Status = jobStatus.Running
		err := cmd.Run()
		handleAPI(err, job, "Job Exited With Error")

		// When finished, mark as stopped
		job.Status = jobStatus.Stopped
	}
}

// Stop a job Task
func (job *taskJob) Stop() {

	var cmd *exec.Cmd

	// Mark Job as stopping, Clear error Message
	job.Status = jobStatus.Stopping

	// Makes sure to run the correct stop sequence
	switch job.Type {
	case "docker":
		cmd = job.DockerStop()
	}

	err := cmd.Run()
	handleAPI(err, job, "Failed to Stop Job")
}

// Stop sequence for a docker job
func (job *taskJob) DockerStop() *exec.Cmd {
	composeLoc := folder.Docker + job.ID + "/docker-compose.yml"

	// Stopping the container
	cmd := exec.Command("docker-compose", "-f", composeLoc, "down")

	return cmd
}

// Get cmd for a Docker Job
func (job *taskJob) DockerCmd() *exec.Cmd {

	// All Location Data
	id := job.ID
	url := job.URL
	logsLoc := folder.Logs + id
	composeLoc := folder.Docker + id + "/docker-compose.yml"

	// Make url, read the compose file
	resp, err := http.Get(url)
	handleAPI(err, job, "Failed to get compose URL")
	defer resp.Body.Close()
	file, err := ioutil.ReadAll(resp.Body)
	handleAPI(err, job, "Failed to read compose file")

	// Make directory for docker and logs and save file
	os.Mkdir(folder.Docker+id, 0777)
	os.Mkdir(logsLoc, 0777)
	err = ioutil.WriteFile(composeLoc, file, 0777)
	handleAPI(err, job, "Failed to write to file")


	// Remove any previous containers
	out, err := exec.Command("docker-compose", "-f", composeLoc, "down").CombinedOutput()
	handleAPI(err, job, string(out))

	// Make sure Config Is valid
	out, err = exec.Command("docker-compose", "-f", composeLoc, "config").CombinedOutput()
	handleAPI(err, job, string(out))

	// pull lastest images
	out, err = exec.Command("docker-compose", "-f", composeLoc, "pull").CombinedOutput()
	handleAPI(err, job, string(out))

	// Get logging Running
	cmd := exec.Command("docker-compose", "-f", composeLoc, "up", "--no-color")

	outfile, err := os.Create(logsLoc + "/info.log")
	handleAPI(err, job, "Failed to create log file")
	cmd.Env = job.insertEnviroment()
	cmd.Stdout = outfile

	return cmd
}

func (job *taskJob) insertEnviroment() []string {
	var environ []string

	for key, val := range job.Enviroment {
		str := key + "=" + val
		environ = append(environ, str)
	}

	return environ
}
