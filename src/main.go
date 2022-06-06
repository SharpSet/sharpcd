package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var secretFlag string
var remoteFile string
var sharpURL string

// Create Flags needed
func init() {
	flag.StringVar(&secretFlag, "secret", "", "Put secret as a arg for automation tasks")
	flag.StringVar(&remoteFile, "remotefile", "", "Location of Remote sharpcd.yml file")
	flag.StringVar(&sharpURL, "sharpurl", "", "Location of SharpCD server (Will override sharpcd.yml)")

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

	- sharpcd trak alljobs {location}
		Get info on all jobs

	- sharpcd trak job {location} {job_id}
		Get info on job with logging

	- sharpcd trak list {location}
		Get all jobs running on sharpcd server

	- sharpcd trak logs {location} {job_id}
		Get Logs from a Job

Flags:
		`)

		flag.PrintDefaults()
		fmt.Println()
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
			flag.Usage()
			log.Fatal("This subcommand does not exist!")
		}
	}
	return
}

// Get the local directory
// Method changes depending on enviroment
func getDir() string {
	var exPath string
	var err error

	ex, err := os.Executable()
	handle(err, "Failed to get dir")
	exPath = filepath.Dir(ex)

	return exPath
}
