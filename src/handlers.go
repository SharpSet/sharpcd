package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

// Hander for requests made to /api/
func httpHandleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := statCode.Accepted

	// Do Common Checks
	commonHandlerChecks(r, &status)

	resp := response{}
	w.WriteHeader(status)

	err := getAPIData(r, &resp)
	handleStatus(err, statCode.FailedLogRead, &status)

	// If all of that passed, send message showing success
	if status == statCode.Accepted {

	} else {
		resp.Message = getFailMessage(status)
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

	case statCode.NotPostMethod:
		return "SharpCD: Only accepting POST requests"

	case statCode.CommDoesNotExist:
		return "SharpCD: The task type requested does not exist!"

	case statCode.FailedLogRead:
		return "SharpCD: Could not read log file for job"

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

// Check the correct http method is used
func checkMethod(method string) error {
	if method != "POST" {
		return errors.New("Wrong Method")
	}

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
