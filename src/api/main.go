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

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func getTLSConfig() (*tls.Config, error) {
	crt, err := ioutil.ReadFile("../certs/localhost.crt")
	if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile("../certs/localhost.key")
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

func main() {
	l := log.New(os.Stdout, "saltgram-api-gateway", log.LstdFlags)
	l.Printf("Starting API Gateway on port: %s\n", os.Getenv("SALT_API_PORT"))

	serverMux := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{fmt.Sprintf("https://localhost:%s", os.Getenv("SALT_WEB_PORT"))},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
		Debug:            true,
	})

	tlsConfig, err := getTLSConfig()
	if err != nil {
		l.Fatalf("[ERROR] configuring tls: %v\n", err)
	}

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("SALT_API_PORT")),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      c.Handler(serverMux),
		TLSConfig:    tlsConfig,
	}

	go func() {
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			l.Fatalf("[ERROR] while serving: %v\n", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	sig := <-signalChan
	l.Println("Recieved terminate, graceful shutdown with sigtype:", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(tc)
}
