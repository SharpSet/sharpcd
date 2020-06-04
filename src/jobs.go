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
		URL:  payload.GitURL + payload.Compose}

	// If a job with that ID already exists
	if job := getJob(payload.ID); job != nil {

		// Assign its ID to the old job ID
		newJob.ID = job.ID

		// Stop the old job
		job.Status = jobStatus.Stopping
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
		cmd = job.DockerRun()
	}

	// Mark as Running
	job.Status = jobStatus.Running
	err := cmd.Run()
	handleAPI(err, job, "Failed to Run Job")

	// When finished, mark as stopped
	job.Status = jobStatus.Stopped
}

// Stop a job Task
func (job *taskJob) Stop() {

	var cmd *exec.Cmd

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

// Run sequence for a Docker job
func (job *taskJob) DockerRun() *exec.Cmd {

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
	os.Mkdir(folder.Docker + id, 0777)
	os.Mkdir(logsLoc, 0777)
	err = ioutil.WriteFile(composeLoc, file, 0777)
	handleAPI(err, job, "Failed to write to file")

	job.Status = jobStatus.Building

	// Build Commands
	err = exec.Command("docker-compose", "-f", composeLoc, "down").Run()
	handleAPI(err, job, "Failed to restart docker job")
	err = exec.Command("docker-compose", "-f", composeLoc, "pull").Run()
	handleAPI(err, job, "Failed to Pull New images for docker job")

	// Get logging Running
	cmd := exec.Command("docker-compose", "-f", composeLoc, "up", "--no-color")

	outfile, err := os.Create(logsLoc + "/info.log")
	handleAPI(err, job, "Failed to create log file")
	cmd.Stdout = outfile

	return cmd
}
