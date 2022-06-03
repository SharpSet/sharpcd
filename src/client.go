package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func client() {

	var resp bool
	var con config
	var err error
	var file []byte

	if len(remoteFile) != 0 {
		resp, err := http.Get(remoteFile)
		handle(err, "Failed to download remote sharpcd.yml")
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			file, err = ioutil.ReadAll(resp.Body)
			handle(err, "Failed to read remote sharpcd.yml")
		}
	} else {
		// Get config, get data from it
		file, err = ioutil.ReadFile("./sharpcd.yml")
		handle(err, "Failed to read and extract sharpcd.yml")
	}

	err = yaml.Unmarshal(file, &con)
	handle(err, "Failed to read yaml from sharpcd.yml")

	compareVersions(con.Version)

	// store tasks already run
	var tasksRun []string

	// variable to store nesting level
	var level int

	// POST to sharpcd server for each task
	for id, task := range con.Tasks {

		// Make ID Lower Case
		id := strings.ToLower(id)
		id = strings.ReplaceAll(id, " ", "_")

		resp = runTask(id, task, &tasksRun, con, level)
	}

	if resp {
		fmt.Println("At least one task Failed!")
		os.Exit(1)
	}
}

func runTask(id string, task task, tasksRun *[]string, con config, level int) (response bool) {
	response = false

	// if level is above 10, exit
	if level > 10 {
		fmt.Println("Too many nested tasks, exiting from", id)
		fmt.Println("This is likely a bug in your config. Check to ensure that two tasks are not dependent on each other.")
		fmt.Println(id)
		os.Exit(1)
	}

	payload := postData{
		ID:         id,
		Name:       task.Name,
		Type:       task.Type,
		GitURL:     task.GitURL,
		Command:    task.Command,
		Compose:    task.Compose,
		Enviroment: getEnviroment(task.Envfile),
		Registry:   task.Registry,
		Secret:     getSec(),
		Version:    sharpCDVersion}

	// check for task dependencies
	if len(task.Depends) != 0 {
		for _, taskDep := range task.Depends {
			level++
			runTask(taskDep, con.Tasks[taskDep], tasksRun, con, level)
		}
	}

	alreadyRun := false
	// check that task has not already been run
	for _, taskRun := range *tasksRun {
		if taskRun == payload.ID {
			alreadyRun = true
			break
		}
	}

	if !alreadyRun {

		var url string

		// if the sharpurl flag is set, use it

		// print sharpURL
		if len(sharpURL) != 0 {
			url = sharpURL
		} else {
			url = task.SharpURL
		}

		// Make POST request and let user know if successful
		body, code := post(payload, url)
		if code == statCode.Accepted {
			fmt.Printf("Task [%s] succesfully sent!\n", task.Name)
			fmt.Println("=================")
			err := postCommChecks(task, id, url)
			if err != nil {
				response = true
			}
			fmt.Println("")
			fmt.Println("")
		} else {
			fmt.Println(body.Message)
			fmt.Printf("Task %s Failed!\n\n", task.Name)
			os.Exit(1)
		}

		// Add task to tasksRun
		*tasksRun = append(*tasksRun, payload.ID)

	}

	return response
}
