package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func server() {
	http.HandleFunc("/", command)
	log.Fatal(http.ListenAndServe(":5666", nil))
}

// Takes a command with POST data
func command(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case "POST":
		text := response{
			Name: "Adam McArthur",
			Value: "Succesful Attempt!",
		}
		json.NewEncoder(w).Encode(text)
		return
	default:
		http.Error(w, "Sorry, SharpCD can only take POST requests.", 422)
		return
	}
}
