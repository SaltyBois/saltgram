package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"saltgram/handlers"
	"time"
)

func main() {
	l := log.New(os.Stdout, "saltgram", log.LstdFlags)
	
	usersHandler := handlers.NewUsers(l)

	serverMux := http.NewServeMux()
	serverMux.Handle("/users", usersHandler)

	server := &http.Server{
		Addr: ":8081",
		Handler: serverMux,
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, os.Kill)

	sig := <- signalChan
	l.Println("Recieved terminate, graceful shutdown with sigtype:", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()
	server.Shutdown(tc)
}
