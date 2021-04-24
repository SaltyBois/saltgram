package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	spa "github.com/roberthodgen/spa-server"
)

func getTLSConfig() (*tls.Config, error) {
	crt, err := ioutil.ReadFile("../../certs/localhost.crt")
	if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile("../../certs/localhost.key")
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
	}, nil
}

func hstsMiddleware(h http.Handler) func (http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=86400")
		h.ServeHTTP(w, r)
	}
}

func main() {
	log.Printf("Running web server on port: %v\n", os.Getenv("SALT_WEB_PORT"))
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", hstsMiddleware(spa.SpaHandler("../frontend/dist", "index.html")))

	tlsConfig, err := getTLSConfig()
	if err != nil {
		log.Fatalf("[ERROR] getting TLS config: %v\n", err)
	}

	server := http.Server{
		Addr: fmt.Sprintf(":%s", os.Getenv("SALT_WEB_PORT")),
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
		TLSConfig: tlsConfig,
		Handler: serverMux,
	}

	go func() {
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatalf("[ERROR] while serving: %v\n", err)
		}
	}()


	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	sig := <-signalChan
	log.Println("Recieved terminate, graceful shutdown with sigtype:", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(tc)
}