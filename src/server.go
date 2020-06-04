package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
)

var certLoc = folder.Private + "/server.crt"
var keyLoc = folder.Private + "/server.key"

func server() {
	// Check if keys have been created
	_, err := os.Stat(certLoc)
	_, err = os.Stat(keyLoc)
	handle(err, "Failed to load openssl keys")


	// Handler Functions
	mux := http.NewServeMux()
	mux.HandleFunc("/", httpHandleTask)
	mux.HandleFunc("/api/", httpHandleAPI)

	// Set config for better security
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

	// Run server
	srv := &http.Server{
		Addr:         ":5666",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	fmt.Println("Server Successfully Running!")
	log.Fatal(srv.ListenAndServeTLS(certLoc, keyLoc))
}
