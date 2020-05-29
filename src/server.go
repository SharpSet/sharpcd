package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"
)

func server() {
	// Check if keys have been created
	_, err := os.Stat("private/server.crt")
	_, err = os.Stat("private/server.key")
	clientErrCheck(err, "Failed to load openssl keys")

	mux := http.NewServeMux()
	mux.HandleFunc("/", command)
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	srv := &http.Server{
		Addr:         ":5666",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	fmt.Println("Server Successfully Running!")
	log.Fatal(srv.ListenAndServeTLS("private/server.crt", "private/server.key"))
}

// Takes a command with POST data
func command(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := statusAcceptedTask
	statuspointer := &status

	// Check Method
	err := checkMethod(r.Method)
	serverErrCheck(err, statusNotPostMethod, statuspointer)

	// Check Body
	body, err := ioutil.ReadAll(r.Body)
	serverErrCheck(err, statusFailedToReadBody, statuspointer)

	// Unmarshal json data
	var payload postData
	err = json.Unmarshal(body, &payload)
	serverErrCheck(err, statusBodyNotJSON, statuspointer)

	// Check URL
	err = checkURL(payload.GitURL)
	serverErrCheck(err, statusBannedURL, statuspointer)

	// Check Password
	err = checkPass(payload.Key)
	serverErrCheck(err, statusIncorrectPass, statuspointer)

	w.WriteHeader(status)

	// If all of that passed, send message showing success
	if status == statusAcceptedTask {
		resp := response{}

		json.NewEncoder(w).Encode(resp)
	} else {
		resp := response{
			Message: getFailMessage(status)}
		json.NewEncoder(w).Encode(resp)
	}

	return
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

func getFailMessage(status int) string {
	switch status {
	case statusBannedURL:
		return "SharpCD: This URL is not allowed on this server"

	case statusBodyNotJSON:
		return "SharpCD: The body of the request is not valid JSON"

	case statusFailedToReadBody:
		return "SharpCD: The body of the request could not be read"

	case statusIncorrectPass:
		return "SharpCD: Incorrect Password"

	case statusNotPostMethod:
		return "SharpCD: Only accepting POST requests"

	default:
		return "No Fail Message"
	}
}
