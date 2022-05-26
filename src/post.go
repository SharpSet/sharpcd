package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

// Makes POST Request adn reads response
func post(payload postData, url string) (response, int) {
	// Create POST request with JSON
	jsonStr, err := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	handle(err, "Failed to create request")

	// Create client
	// Allow self-signed certs
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 5 * time.Second,
	}

	// Do Request
	resp, err := client.Do(req)
	handle(err, "Failed to do POST request to "+url)
	defer resp.Body.Close()

	// Read Body and Status
	body, err := ioutil.ReadAll(resp.Body)
	handle(err, "Failed to read body of response")

	var respBody response
	err = json.Unmarshal(body, &respBody)
	handle(err, "Failed to convert repsonse to JSON")

	return respBody, resp.StatusCode
}

// Get the enviromental vars from a env file
func getEnviroment(loc string) map[string]string {

	var env map[string]string

	// If no envfile specified end function
	if len(loc) == 0 {
		return env
	}

	// Get env contents
	env, err := godotenv.Read(loc)
	handle(err, "Failed to Load .env")
	return env
}

func postCommChecks(t task, id string) error {
	jobURL := t.SharpURL + "/api/job/" + id
	logsURL := t.SharpURL + "/api/logs/" + id
	payload := postData{
		Secret: getSec()}
	buildingTriggered := false
	runningTriggered := false
	lastIssue := ""
	counter := 0

	fmt.Println("Waiting on server response...")

	// Ensure task hasn't stopped unexpectantly early
	for {
		resp, code := post(payload, jobURL)
		if code != statCode.Accepted {
			log.Fatal("Something went wrong using the API!")
			return errors.New("Bad API")
		}

		if resp.Job == nil {
			continue
		}
		job := resp.Job

		stopped := job.Status == jobStatus.Stopped && runningTriggered
		errored := job.Status == jobStatus.Errored
		building := job.Status == jobStatus.Building && !buildingTriggered
		running := job.Status == jobStatus.Running && !runningTriggered

		// Marks that a new build has started
		if building {
			buildingTriggered = true
			fmt.Println("The Task is now building a job")
		}

		if running {
			runningTriggered = true
			fmt.Println("Task Has Successfully started running!")
			fmt.Println("Making sure it does not stop unexpectedly...")
		}

		// Marks some sort of error
		if errored || stopped {
			fmt.Println("Task stopped running!")
			fmt.Println("Error Message: " + job.ErrMsg)

			logResp, code := post(payload, logsURL)
			if code != statCode.Accepted {
				log.Fatal("Something went wrong using the API!")
				return errors.New("Bad API")
			}
			file := logResp.Message

			if file != "" {
				fmt.Println("Logs File:")
				fmt.Println(file)
			}

			return errors.New("Bad Task")
		}

		if lastIssue != job.Issue {
			fmt.Println("Non fatal Issue found: " + job.Issue)
			lastIssue = job.Issue
		}

		// If 7 seconds has elapsed, comsider it started properly
		if counter > 7 {
			fmt.Println("Task has started Properly!")
			return nil
		}

		if runningTriggered {
			counter++
		}

		time.Sleep(1 * time.Second)
	}
}
