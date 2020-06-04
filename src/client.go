package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

func client() {

	// Get config, get data from it
	f, err := ioutil.ReadFile("./sharpcd.yml")
	var con config
	err = yaml.Unmarshal(f, &con)
	handle(err, "Failed to read and extract sharpcd.yml")

	// POST to sharpcd server for each task
	for id, task := range con.Tasks {

		// Make ID Lower Case
		id := strings.ToLower(id)
		id = strings.ReplaceAll(id, " ", "_")

		payload := postData{
			ID:         id,
			Name:       task.Name,
			Type:       task.Type,
			GitURL:     task.GitURL,
			Command:    task.Command,
			Compose:    task.Compose,
			Enviroment: getEnviroment(task.Envfile),
			Secret:     getSec()}

		// Make POST request and let user know if successful
		err = post(payload, task.SharpURL)
		if err == nil {
			fmt.Printf("Task %s succesfully sent!\n\n", task.Name)
		} else {
			fmt.Println(err)
			fmt.Printf("Task %s Failed!\n\n", task.Name)
		}
	}
}

// Makes POST Request adn reads response
func post(payload postData, url string) error {
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
	}

	// Do Request
	resp, err := client.Do(req)
	handle(err, "Failed to do POST request")
	defer resp.Body.Close()

	// Read Body and Status
	body, err := ioutil.ReadAll(resp.Body)
	handle(err, "Failed to read body of response")

	var respBody response
	err = json.Unmarshal(body, &respBody)
	handle(err, "Failed to convert repsonse to JSON")

	// Checks if status is OK
	if resp.StatusCode != statCode.Accepted {
		return errors.New(respBody.Message)
	}

	return nil
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
