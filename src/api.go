package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func handleAPI(w http.ResponseWriter, r *http.Request) {
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

	// Check Password
	err = checkPass(payload.Key)
	checkStatus(err, statusIncorrectPass, statuspointer)

	w.WriteHeader(status)

	// If all of that passed, send message showing success
	if status == statusAcceptedTask {
		resp := apiResponse{
			Procs: allProcs.List}

		json.NewEncoder(w).Encode(resp)

	} else {
		resp := response{
			Message: getFailMessage(status)}
		json.NewEncoder(w).Encode(resp)
	}

	return
}
