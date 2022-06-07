package main

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

var hashLoc = folder.Private + "/hash.secret"
var secCache string
var sharpURLCache string

func checkSecret(sec string) error {

	// Get hash from file
	hash, err := ioutil.ReadFile(hashLoc)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(sec))
	return err
}

func getSec() string {

	// if cache is set, return that
	if len(secCache) != 0 {
		return secCache
	}

	// If secret is not in args
	if len(secretFlag) == 0 {

		// Get secret from user
		fmt.Println("Enter secret: ")
		sec, err := terminal.ReadPassword(0)
		handle(err, "Failed to read secret")

		// set cache
		secCache = string(sec)
		return string(sec)
	}

	// Otherwise, return flag with secret
	return secretFlag
}

// Sets secret into file
func setSec() {

	// Get secret from user
	sec := getSec()

	// Use Bcyrpt to hash and salt
	hash, err := bcrypt.GenerateFromPassword([]byte(sec), bcrypt.MinCost)
	handle(err, "Failed to generate secret hash")

	// Store file
	err = ioutil.WriteFile(hashLoc, hash, 0644)
	handle(err, "Failed to save hash")

	if err == nil {
		fmt.Println("Secret successfully created!")
	}
}

func getSharpURL() string {

	// if cache is set, return that
	if len(sharpURLCache) != 0 {
		return sharpURLCache
	}

	// If secret is not in args
	if len(sharpURL) == 0 {

		// Get sharpurl from user
		fmt.Println()
		fmt.Println("Enter URl to send request to (e.g https://localhost:5666): ")
		var url string
		_, err := fmt.Scanln(&url)
		handle(err, "Failed to read url")
		fmt.Println()

		// set cache
		sharpURLCache = string(url)
		return string(url)
	}

	// Otherwise, return flag with secret
	return sharpURL
}
