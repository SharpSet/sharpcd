package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"
)

// Checks for client err
// Records as Fatal
func check(e error, msg string) {
	// Try and get SHARPDEV var
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cannot read enviroment")
	}

	if e != nil {
		if os.Getenv("SHARPDEV") == "TRUE" {
			fmt.Println(e)
		}
		log.Fatal(msg)
	}
}

// checks for server err
// Writes response given to header
func checkStatus(e error, status int, passedChecks *int) {
	if e != nil {
		*passedChecks = status
	}
}

// checks for server err
// Writes response to API call
func checkAPI(e error, proc *taskProcess) {
	if e != nil {
		proc.Err = e
		proc.Status = procStatError
	}
}

func checkMethod(method string) error {
	if method != "POST" {
		return errors.New("Wrong Method")
	}

	return nil
}

func checkPass(pwd string) error {

	// Get hash from file
	hash, err := ioutil.ReadFile("./private/hash.key")
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(pwd))
	return err
}

// Checks if URLs are okay
func checkURL(taskURL string) error {

	// Read filter file and extract array of allowed urls
	file, err := ioutil.ReadFile("./data/filter.yml")
	if err != nil {
		return err
	}

	// Load into YAML struct
	var f filter
	err = yaml.Unmarshal(file, &f)
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

func checkTaskType(payload postData) error {
	switch payload.Type {
	case "docker":
		return nil
	default:
		return errors.New("Type Doesn't exist")
	}
}
