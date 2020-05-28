package main

import (
	"os"
	"log"
	"net/http"
)

func main() {
	if len(os.Args) == 1 {
		client()
	} else {
		var arg = os.Args[1]

		switch arg {
			case "server":
				server()
			default:
				log.Fatal("This subcommand does not exist!")
		}
	}
	return
}


func clientErrCheck(e error, msg string) {
	if e != nil {
		log.Fatal(msg)
	}
}

func serverErrCheck(w http.ResponseWriter, e error, status int) {
	if e != nil {
		w.WriteHeader(status)
	}
}
