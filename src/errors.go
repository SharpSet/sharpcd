package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Checks for client err
// Records as Fatal
func handle(e error, msg string) {
	// Try and get SHARPDEV var
	godotenv.Load()

	if e != nil {
		if os.Getenv("SHARPDEV") == "TRUE" {
			fmt.Println(e)
		}
		log.Fatal(msg)
	}
}

// checks for server err
// Writes response given to header
func handleStatus(e error, status int, passedChecks *int) {
	if e != nil {
		*passedChecks = status
	}
}

// checks for server err
// Writes response to API call
func handleAPI(e error, job *taskJob, msg string) {
	if e != nil && job.ErrMsg == "" {
		job.ErrMsg = msg
		job.Status = jobStatus.Errored
	}
}
