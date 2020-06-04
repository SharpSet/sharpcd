package main

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

var hashLoc = folder.Private + "/hash.secret"

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

	// If secret is not in args
	if len(secretFlag) == 0 {

		// Get secret from user
		fmt.Println("Enter secret: ")
		sec, err := terminal.ReadPassword(0)
		handle(err, "Failed to read secret")
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
