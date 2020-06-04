package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

var passFlag string

// Create Flags needed
func init() {
	flag.StringVar(&passFlag, "pass", "", "Put password as a arg for automation tasks")

	// Creates Helper Function
	flag.Usage = func() {
		fmt.Println(`
Args of SharpCD:

	- server: Run the sharpcd server
	- setpass: Set the password for API and Task Calls

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
		case "setpass":
			setPwd()
		default:
			log.Fatal("This subcommand does not exist!")
		}
	}
	return
}

func getPwd() string {

	// If password is not in args
	if len(passFlag) == 0 {

		// Get password from user
		fmt.Println("Enter password: ")
		pwd, err := terminal.ReadPassword(0)
		check(err, "Failed to read password")
		return string(pwd)
	}

	// Otherwise, return flag with password
	return passFlag
}

// Sets password into file
func setPwd() {

	// Get password from user
	pwd := getPwd()

	// Use Bcyrpt to hash and salt
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	check(err, "Failed to generate password hash")

	// Store file
	err = ioutil.WriteFile("./private/hash.key", hash, 0644)
	check(err, "Failed to save hash")

	if err == nil {
		fmt.Println("Password successfully created!")
	}
}
