package main

import (
	"os"
)

func main() {
	var arg = os.Args[1]

	switch arg {
	case "server":
		server()
	}
	return
}
