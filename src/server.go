package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	log.Fatal(srv.ListenAndServeTLS("private/server.crt", "private/server.key"))
}

// Takes a command with POST data
func command(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case "POST":

		body, err := ioutil.ReadAll(r.Body)
		serverErrCheck(w, err, http.StatusBadRequest)

		// Unmarshal json data
		var payload postData
		err = json.Unmarshal(body, &payload)
		serverErrCheck(w, err, http.StatusBadRequest)

		if err == nil {
			w.WriteHeader(http.StatusOK)
		}

		json.NewEncoder(w).Encode(payload)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
