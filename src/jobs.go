package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
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
		Name:       payload.Name,
		Type:       payload.Type,
		URL:        payload.GitURL + payload.Compose,
		Enviroment: payload.Enviroment,
		Registry:   payload.Registry,
		Reconnect:  false}

	// If a job with that ID already exists
	if job := getJob(payload.ID); job != nil {

		// Assign its ID to the old job ID
		newJob.ID = job.ID

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

		job.Status = jobStatus.Stopped

	}
}

// Get cmd for a Docker Job
func (job *taskJob) DockerCmd() *exec.Cmd {

	// All Location Data
	id := job.ID
	url := job.URL
	logsLoc := folder.Docker + id
	composeLoc := folder.Docker + id + "/docker-compose.yml"

	var err error

	if job.Reconnect != true {
		// Get github token
		f, err := readFilter()
		if err != nil {
			handleAPI(err, job, "Failed to get Token")
		}

		// Make url, read the compose file
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			handleAPI(err, job, "Failed to build request")
		}

		if f.Token != "" {
			req.Header.Set("Authorization", "token "+f.Token)
			req.Header.Set("Accept", "application/vnd.github.v3.raw")
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			handleAPI(err, job, "Failed to get compose URL")
		}
		defer resp.Body.Close()
		file, err := ioutil.ReadAll(resp.Body)
		handleAPI(err, job, "Failed to read compose file")

		// Make directory for docker and logs and save file
		os.Mkdir(folder.Docker+id, 0777)
		os.Mkdir(logsLoc, 0777)

		if strings.Contains(string(file), "404") {
			err = errors.New("404 in Compose File")
			text := "Github Token invalid or wrong compose URL"
			handleAPI(err, job, text)
			job.Issue = text
		}

		// If Token was valid and compose file not empty
		if err == nil {

			// Ensure ComposeLoc is Empty
			_, err = os.Stat(composeLoc)
			if err == nil {
				err = os.Remove(composeLoc)

				if err != nil {
					handleAPI(err, job, "Failed to Remove")
				}
			}

			// Write to file
			err = ioutil.WriteFile(composeLoc, file, 0777)
			handleAPI(err, job, "Failed to write to file")
		}

		if job.Registry != "" {
			job.dockerLogin()
		}

		// Pull new images
		job.buildCommand("-f", composeLoc, "pull")

		// Make sure config is valid
		err = job.buildCommand("-f", composeLoc, "up", "--no-start")
		if err == nil {
			// Remove any previous containers
			// Deals with any network active endpoints
			job.buildCommand("-f", composeLoc, "down", "--remove-orphans")

			// Run Code
			job.buildCommand("-f", composeLoc, "up", "-d")

			if job.Registry != "" {
				job.dockerLogout()
			}
		} else {
			return nil
		}
	}

	// Get logging Running
	cmd := exec.Command("docker-compose", "-f", composeLoc, "logs", "-f", "--no-color")

	outfile, err := os.Create(logsLoc + "/info.log")
	handleAPI(err, job, "Failed to create log file")
	cmd.Env = job.insertEnviroment()
	cmd.Stdout = outfile

	return cmd
}

func (job *taskJob) insertEnviroment() []string {

	var err error
	envfile := folder.Docker + job.ID + "/.env"
	var environ []string

	// If there is enviroment vars to use
	if len(job.Enviroment) != 0 {

		// Apply them
		_, err = os.Create(envfile)
		err = godotenv.Write(job.Enviroment, envfile)

		handleAPI(err, job, "Failed to write to Envfile")
	} else {
		// Check if file exists for .env vars

		_, err = os.Stat(envfile)
		if err == nil {
			job.Issue = "EnvFile exists but no Job Env Vars was given"
		}
	}

	fileEnviroment, err := godotenv.Read(envfile)
	if err == nil {
		for key, val := range fileEnviroment {
			str := key + "=" + val
			environ = append(environ, str)
		}

		return append(os.Environ(), environ...)
	}

	return nil
}

func (job *taskJob) dockerLogin() {
	cmd := exec.Command("docker", "login", "-u", job.Enviroment["DOCKER_USER"], "-p", job.Enviroment["DOCKER_PASS"], job.Registry)
	out, err := cmd.CombinedOutput()
	errMsg := string(out)
	handleAPI(err, job, errMsg)
}

func (job *taskJob) dockerLogout() {
	cmd := exec.Command("docker", "logout", job.Registry)
	out, err := cmd.CombinedOutput()
	errMsg := string(out)
	handleAPI(err, job, errMsg)
}

func (job *taskJob) buildCommand(args ...string) error {
	var errMsg string
	cmd := exec.Command("docker-compose", args...)
	cmd.Env = job.insertEnviroment()
	out, err := cmd.CombinedOutput()

	// Add conditions for volumes and networks
	if strings.Contains(string(out), "404") {
		errMsg = "No Compose File Found!"
		handleAPI(err, job, errMsg)
		return err
	} else if strings.Contains(string(out), "manually using `") {

		for {
			cmd := exec.Command("docker-compose", args...)
			out, err = cmd.CombinedOutput()

			if strings.Contains(string(out), "manually using `") {
				// Find Docker Command
				re := regexp.MustCompile("`(.*)`")
				command := strings.ReplaceAll(string(re.Find(out)), "`", "")
				commands := strings.Split(command, " ")

				// Create Missing Element
				cmd := exec.Command(commands[0], commands[1:]...)

				// Handle Errors
				out, err := cmd.CombinedOutput()
				handleAPI(err, job, string(out))
			} else {
				break
			}
		}

	} else {
		errMsg = string(out)
		handleAPI(err, job, errMsg)
		return err
	}

	return nil

}
