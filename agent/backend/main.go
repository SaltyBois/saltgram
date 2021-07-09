package main

import (
	"agent/data"
	"agent/handlers"
	"agent/internal"
	"agent/pki"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/cors"
)

func main() {
	l := internal.NewLogger("saltgram-agent")
	pkiHandler := pki.Init()
	cert, err := pkiHandler.RegisterSaltgramService("saltgram-agent")
	pkiHandler.RegisterSaltgramService("agent-web")
	if err != nil {
		l.L.Fatalf("failed to register service: %v", err)
	}
	s := internal.NewService(l.L)
	err = s.Init("saltgram-agent", cert.CertPEM, cert.PrivateKeyPEM, pkiHandler.RootCA.CertPEM)

	if err != nil {
		l.L.Fatalf("failed to init agent: %v\n", err)
	}

	db := data.NewDBConn(l.L)
	err = db.ConnectToDb()
	db.MigradeData()
	if err != nil {
		l.L.Fatalf("failed to connect to agent db: %v\n", err)
	}

	agentHandler := handlers.NewAgent(l.L, db)
	agentRouter := s.S.PathPrefix("/agent").Subrouter()
	agentRouter.HandleFunc("/signup", agentHandler.Signup).Methods(http.MethodPost)
	agentRouter.HandleFunc("/signin", agentHandler.SignIn).Methods(http.MethodPost)
	agentRouter.HandleFunc("/isagent", agentHandler.IsAgent).Methods(http.MethodGet)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{fmt.Sprintf("https://localhost:%s", os.Getenv("AGENT_WEB_PORT"))},
		//AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
		Debug:            true,
	})

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("AGENT_PORT")),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      c.Handler(s.S),
		TLSConfig:    s.TLS.TC,
	}

	go func() {
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			l.L.Fatalf("while serving: %v\n", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	sig := <-signalChan
	l.L.Printf("Recieved terminate, graceful shutdown with sigtype: %v\n", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(tc)
}
