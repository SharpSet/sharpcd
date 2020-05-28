package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
	"crypto/tls"
)


func client() {

	// Get config, get data from it
	f, err := ioutil.ReadFile("./config.yml")
	var con config
	err = yaml.Unmarshal(f, &con)
	clientErrCheck(err, "Failed to read and extract config.yml")

	// testing making enviroment map
	env := make(map[string]string)
	env["TEST"] = "HELLO"

	for _, task := range con.Tasks {
		payload := postData{
			Name: task.Name,
			Type: task.Type,
			GitURL: task.GitURL,
			Command: task.Command,
			Enviroment: env,
			Key: "password"}

		post(payload, task.SharpURL)
	}
}

func post (payload postData, url string) {
    jsonStr, err := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	clientErrCheck(err, "Failed to create request")

    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
    resp, err := client.Do(req)
	clientErrCheck(err, "Failed to do POST request")
    defer resp.Body.Close()

	fmt.Println("response Status:", resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	clientErrCheck(err, "Failed to read body of response")
    fmt.Println("response Body:", string(body))
}
