package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Checks which API section is required and
// runs the appropreiate command
func getAPIData(r *http.Request, resp *response) error {
	var err error

	path := strings.Split(r.URL.Path[5:], "/")
	switch path[0] {
	case "jobs":
		for _, job := range allJobs {
			resp.Jobs = append(resp.Jobs, job)
		}

		sort.Slice(resp.Jobs, func(i, j int) bool {
			return resp.Jobs[i].ID < resp.Jobs[j].ID
		})
		return nil
	case "job":
		resp.Job, err = getJobs(path[1])
		return err
	case "logs":
		resp.Message, err = getLogs(path[1])
		return err
	case "logsfeed":
		resp.Message, err = getLogsFeed(path[1])
		return err
	}

	return nil
}

// Gets the logs from the task ID's file
func getLogs(path string) (string, error) {
	logs := folder.Docker + path + "/info.log"
	file, err := os.Open(logs)
	if err != nil {
		return "", err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func getLogsFeed(path string) (string, error) {
	logs := folder.Docker + path + "/info.log"

	cmd := exec.Command("tail", "-11", logs)
	out, err := cmd.CombinedOutput()
	msg := string(out)
	if err != nil {
		return msg, errors.New("Failed to run Docker Compose")
	}

	return msg, nil
}

func getJobs(path string) (*taskJob, error) {
	var emptyJob *taskJob

	for id, job := range allJobs {
		if id == path {
			err := checkJobStatus(job)
			return job, err
		}
	}

	return emptyJob, errors.New("job not found")
}

func checkJobStatus(job *taskJob) error {
	logs, err := getLogs(job.ID)

	exited := strings.Contains(logs, "exited with code")
	if exited && (job.Status != jobStatus.Building) {
		job.Status = jobStatus.Errored
		job.ErrMsg = "A Container has maybe exited unexpectedly"
	}

	return err
}
