package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var secretFlag string

// Create Flags needed
func init() {
	flag.StringVar(&secretFlag, "secret", "", "Put secret as a arg for automation tasks")

	// Creates Helper Function
	flag.Usage = func() {
		fmt.Println(`

Args of SharpCD:

	- server: Run the sharpcd server
	- setsecret: Set the secret for API and Task Calls
	- addfilter: Add a url for a compose file
	- changetoken: Add a token for private github repos
	- removefilter: Remove a url for a compose file
	- version: Returns the Current Version
	- trak: Run the Trak program

Sub Command Trak:

	- alljobs {type}: Get info on all jobs
	- job {type} {id}: Get info on job with logging
	- list {type}: Get all jobs running on sharpcd server

Flags:
		`)

		flag.PrintDefaults()
	}
}

func main() {
	// Parses flags and removes them from args
	flag.Parse()

	if len(flag.Args()) == 0 {
		client()
	} else {
		var arg1 = flag.Args()[0]

		// Subcommands
		switch arg1 {
		case "server":
			server()
		case "trak":
			trak()
		case "help":
			flag.Usage()
		case "setsecret":
			setSec()
		case "addfilter":
			addFilter()
		case "removefilter":
			removeFilter()
		case "changetoken":
			changeToken()
		case "version":
			fmt.Println("Version: " + sharpCDVersion)
		default:
			log.Fatal("This subcommand does not exist!")
			flag.Usage()
		}
	}
	return
}

// Get the local directory
// Method changes depending on enviroment
func getDir() string {
	var exPath string
	var err error

	if os.Getenv("DEV") == "TRUE" {
		exPath, err = os.Getwd()
		handle(err, "Failed to get dir")

	} else {
		ex, err := os.Executable()
		handle(err, "Failed to get dir")
		exPath = filepath.Dir(ex)

	}

	return exPath
}
