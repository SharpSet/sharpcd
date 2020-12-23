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
	- removefilter: Remove a url for a compose file

This will read the sharpdev.yml file
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
		case "setsecret":
			setSec()
		case "addfilter":
			addFilter()
		case "removefilter":
			removeFilter()
		default:
			log.Fatal("This subcommand does not exist!")
		}
	}
	return
}

func getDir() string {
	ex, err := os.Executable()
	handle(err, "Failed to get dir")
	exPath := filepath.Dir(ex)
	return exPath
}
