package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	homeDir := "/"
	if h := os.Getenv("HOME"); h != "" {
		homeDir = h
	}
	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache(filepath.Join(homeDir, ".cache", "golang-autocert")),
		HostPolicy: func(ctx context.Context, host string) error {
			log.Println("HostPolicy:", host)
			return nil
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		fmt.Fprintf(w, "Hello Secure World")
	})

	server := &http.Server{
		Addr:    ":443",
		Handler: mux,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	fmt.Println("listening on " + server.Addr)
	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	server.ListenAndServeTLS("", "")
}
