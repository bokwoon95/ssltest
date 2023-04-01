package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/caddyserver/certmagic"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	certmagic.DefaultACME.Agreed = true
	certConfig := certmagic.NewDefault()
	err := certConfig.ManageAsync(context.Background(), []string{"notebrew.com"})
	if err != nil {
		log.Fatal(err)
	}
	err = certConfig.ObtainCertSync(context.Background(), "notebrew.com")
	if err != nil {
		log.Fatal(err)
	}
	tlsConfig := certConfig.TLSConfig()
	tlsConfig.NextProtos = []string{"h2", "http/1.1", "acme-tls/1"}

	server := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		}),
	}
	fmt.Println("listening on " + server.Addr)
	server.ListenAndServeTLS("", "")
}
