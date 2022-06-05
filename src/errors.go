package main

import (
	"fmt"
	"os"

	ui "github.com/gizak/termui/v3"
	"github.com/joho/godotenv"
)

// Checks for client err
// Records as Fatal
func handle(e error, msg string) {
	// Try and get DEV var
	godotenv.Load()

	if e != nil {
		if os.Getenv("DEV") == "TRUE" {
			fmt.Println(e)
		}
		ui.Close()
		fmt.Println(msg)
		os.Exit(1)
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

	fmt.Println("DEBUG:", job.ErrMsg, job.Status, e)
}
