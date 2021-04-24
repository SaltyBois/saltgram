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
	"saltgram/email/handlers"
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
	l := log.New(os.Stdout, "saltgram-email", log.LstdFlags)

	emailHandler := handlers.NewEmail(l)

	serverMux := mux.NewRouter()

	activationRouter := serverMux.PathPrefix("/activate").Subrouter()
	activationRouter.HandleFunc("", emailHandler.SendActivation).Methods(http.MethodPost)
	activationRouter.HandleFunc("/{token:[A-Za-z0-9]+}", emailHandler.Activate).Methods(http.MethodPut)

	changeRouter := serverMux.PathPrefix("/change").Subrouter()
	changeRouter.HandleFunc("/{token:[A-Za-z0-9]+}", emailHandler.ConfirmReset).Methods(http.MethodPut)
	changeRouter.HandleFunc("", emailHandler.ChangePassword).Methods(http.MethodPost)
	changeRouter.HandleFunc("/forgot", emailHandler.RequestReset).Methods(http.MethodPost)
	
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{fmt.Sprintf("https://localhost:%s", os.Getenv("SALT_API_PORT"))},
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
		Addr: fmt.Sprintf(":%s", os.Getenv("SALT_EMAIL_PORT")),
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout: 120 * time.Second,
		Handler: c.Handler(serverMux),
		TLSConfig: tlsConfig,
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