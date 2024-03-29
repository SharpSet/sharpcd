package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/hashicorp/go-version"
)

// Hander for requests made to /api/
func httpHandleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := statCode.Accepted

	// Do Common Checks
	commonHandlerChecks(r, &status)

	// Check if job exists
	checkJobExists(r, &status)

	resp := response{}
	w.WriteHeader(status)

	err := getAPIData(r, &resp)
	handleStatus(err, statCode.FailedLogRead, &status)

	// If all of that passed, send message showing success
	if status == statCode.Accepted {

	} else {
		resp.Message = getFailMessage(status) + "\nMessage: " + resp.Message
	}

	json.NewEncoder(w).Encode(resp)
	return
}

// handler for commands to /
func httpHandleTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := statCode.Accepted

	payload := commonHandlerChecks(r, &status)

	// Check URL
	err := checkURL(payload.GitURL)
	handleStatus(err, statCode.BannedURL, &status)

	// Check Version
	err = checkVersion(payload.Version)
	handleStatus(err, statCode.WrongVersion, &status)

	// Get task Type
	err = checkTaskType(payload)
	handleStatus(err, statCode.CommDoesNotExist, &status)

	w.WriteHeader(status)
	resp := response{}

	// If all of that passed, send message showing success
	if status == statCode.Accepted {

		// Start Creating Job
		go createJob(payload)

	} else {
		resp.Message = getFailMessage(status)
	}

	json.NewEncoder(w).Encode(resp)
	return
}

// Gets the correct message for each status code
func getFailMessage(status int) string {
	switch status {
	case statCode.BannedURL:
		return "SharpCD: This URL is not allowed on this server"

	case statCode.BodyNotJSON:
		return "SharpCD: The body of the request is not valid JSON"

	case statCode.FailedToReadBody:
		return "SharpCD: The body of the request could not be read"

	case statCode.IncorrectSecret:
		return "SharpCD: Incorrect Secret"

	case statCode.WrongVersion:
		return "SharpCD: Wrong Client Version, expected " + sharpCDVersion + " or above."

	case statCode.NotPostMethod:
		return "SharpCD: Only accepting POST requests"

	case statCode.CommDoesNotExist:
		return "SharpCD: The task type requested does not exist!"

	case statCode.FailedLogRead:
		return "SharpCD: Could not read log file for job"

	case statCode.JobDoesNotExist:
		return "SharpCD: Job does not exist"

	default:
		return "No Fail Message"
	}
}

// Does the common checks for all event handlers
func commonHandlerChecks(r *http.Request, status *int) postData {
	// Check Method
	err := checkMethod(r.Method)
	handleStatus(err, statCode.NotPostMethod, status)

	// Check Body
	body, err := ioutil.ReadAll(r.Body)
	handleStatus(err, statCode.FailedToReadBody, status)

	// Unmarshal json data
	var payload postData
	err = json.Unmarshal(body, &payload)
	handleStatus(err, statCode.BodyNotJSON, status)

	// Check Secret
	err = checkSecret(payload.Secret)
	handleStatus(err, statCode.IncorrectSecret, status)

	return payload
}

func checkJobExists(r *http.Request, status *int) {

	path := strings.Split(r.URL.Path[5:], "/")
	jobID := path[1]

	if path[0] == "job" {
		for _, job := range allJobs {
			if job.ID == jobID {
				return
			}
		}

		*status = statCode.JobDoesNotExist
	}
}

// Check the correct http method is used
func checkMethod(method string) error {
	if method != "POST" {
		return errors.New("Wrong Method")
	}

	return nil
}

// Check that client versions match
func checkVersion(clientVersion string) error {
	if clientVersion != "" {
		v1, err := version.NewVersion(clientVersion)
		v2, err := version.NewVersion(sharpCDVersion)

		if v1.LessThan(v2) {
			err = errors.New("Wrong Client Version")
		}

		return err
	}

	// Means its a non-sharpcd task
	return nil
}

// Checks if URLs are okay
func checkURL(taskURL string) error {

	// Read filter file and extract array of allowed urls
	f, err := readFilter()
	if err != nil {
		return err
	}

	// Parse Task into host and path
	task, err := url.Parse(taskURL)
	if err != nil {
		return err
	}
	taskPath := path.Base(task.Path)

	var foundMatch bool

	// For every allowed url
	for _, allowedURL := range f.Allowed {
		var allowed *url.URL

		allowed, err = url.Parse(allowedURL)
		if err != nil {
			return err
		}
		allowedPath := path.Base(allowed.Path)

		// If they match, mark as such
		if allowed.Host+allowedPath == task.Host+taskPath {
			foundMatch = true
		}
	}

	// If task url does not pass the filter
	if !foundMatch {
		err = errors.New("filter: URL is not allowed")
	}

	return err
}

// Check that the task type exists
func checkTaskType(payload postData) error {
	switch payload.Type {
	case "docker":
		return nil
	default:
		return errors.New("Type Doesn't exist")
	}
}
