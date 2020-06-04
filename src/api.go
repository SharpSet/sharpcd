package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Checks which API section is required and
// runs the appropreiate command
func getAPIData(r *http.Request, resp *response) error {
	var err error

	path := strings.Split(r.URL.Path[5:], "/")
	switch path[0] {
	case "jobs":
		resp.Jobs = allJobs.List
		return nil
	case "logs":
		resp.Message, err = getLogs(path[1])
		return err
	}

	return nil
}

// Gets the logs from the task ID's file
func getLogs(path string) (string, error) {
	logs := folder.Logs + path + "/info.log"
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
