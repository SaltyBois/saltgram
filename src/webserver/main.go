package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"saltgram/internal"
	"time"

	spa "github.com/roberthodgen/spa-server"
)

func hstsMiddleware(h http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=86400")
		h.ServeHTTP(w, r)
	}
}

func main() {
	l := log.New(os.Stdout, "saltgram-web", log.LstdFlags)
	l.Printf("Running web server on port: %v\n", os.Getenv("SALT_WEB_PORT"))
	s := internal.NewService(l)
	s.S.HandleFunc("/", hstsMiddleware(spa.SpaHandler("../frontend/dist", "index.html")))

	err := s.TLS.Init("../../certs/localhost.crt", "../../certs/localhost.key", "")
	if err != nil {
		log.Fatalf("[ERROR] getting TLS config: %v\n", err)
	}

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("SALT_WEB_PORT")),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		TLSConfig:    s.TLS.TC,
		Handler:      s.S,
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
