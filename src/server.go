package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
)

func server() {
	// Check if keys have been created
	_, err := os.Stat("private/server.crt")
	_, err = os.Stat("private/server.key")
	check(err, "Failed to load openssl keys")

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleTask)
	mux.HandleFunc("/api", handleAPI)
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
