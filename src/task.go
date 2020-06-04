package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Takes a task with POST data
func handleTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := statusAcceptedTask
	statuspointer := &status

	// Check Method
	err := checkMethod(r.Method)
	checkStatus(err, statusNotPostMethod, statuspointer)

	// Check Body
	body, err := ioutil.ReadAll(r.Body)
	checkStatus(err, statusFailedToReadBody, statuspointer)

	// Unmarshal json data
	var payload postData
	err = json.Unmarshal(body, &payload)
	checkStatus(err, statusBodyNotJSON, statuspointer)

	// Check URL
	err = checkURL(payload.GitURL)
	checkStatus(err, statusBannedURL, statuspointer)

	// Check Password
	err = checkPass(payload.Key)
	checkStatus(err, statusIncorrectPass, statuspointer)

	// Get task Type
	err = checkTaskType(payload)
	checkStatus(err, statusCommDoesNotExist, statuspointer)

	w.WriteHeader(status)

	// If all of that passed, send message showing success
	if status == statusAcceptedTask {
		resp := response{}

		json.NewEncoder(w).Encode(resp)

		// Start Creating Task
		proc := createProc(payload)
		go proc.Run()

	} else {
		resp := response{
			Message: getFailMessage(status)}
		json.NewEncoder(w).Encode(resp)
	}

	return
}

func getFailMessage(status int) string {
	switch status {
	case statusBannedURL:
		return "SharpCD: This URL is not allowed on this server"

	case statusBodyNotJSON:
		return "SharpCD: The body of the request is not valid JSON"

	case statusFailedToReadBody:
		return "SharpCD: The body of the request could not be read"

	case statusIncorrectPass:
		return "SharpCD: Incorrect Password"

	case statusNotPostMethod:
		return "SharpCD: Only accepting POST requests"

	case statusCommDoesNotExist:
		return "SharpCD: The task type requested does not exist!"

	default:
		return "No Fail Message"
	}
}
