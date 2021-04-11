package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"saltgram/handlers"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "saltgram", log.LstdFlags)

	usersHandler := handlers.NewUsers(l)
	loginHandler := handlers.NewLogin(l)

	serverMux := mux.NewRouter()

	getRouter := serverMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", usersHandler.GetAll)
	getRouter.HandleFunc("/users/{id:[0-9]+}", usersHandler.GetByID)

	postRouter := serverMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/users", usersHandler.Create)
	// TODO(Jovan): Fix later
	postRouter.Use(usersHandler.MiddlewareValidateUser) 

	loginRouter := serverMux.PathPrefix("/login").Subrouter()
	loginRouter.HandleFunc("", loginHandler.Login).Methods(http.MethodPost)
	loginRouter.Use(loginHandler.MiddlewareValidateToken)

	putRouter := serverMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/users/{id:[0-9]+}", usersHandler.Update)
	putRouter.Use(usersHandler.MiddlewareValidateUser)

	// NOTE(Jovan): CORS
	headersOk := gohandlers.AllowedHeaders([]string{"*"})
	originsOk := gohandlers.AllowedOrigins([]string{"*"})
	methodsOk := gohandlers.AllowedMethods([]string{"*"})
	corsHandler := gohandlers.CORS(headersOk, originsOk, methodsOk)

	// h := cors.Default().Handler(serverMux) Works for some reason

	server := &http.Server{
		Addr:         ":8081",
		Handler:      corsHandler(serverMux),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
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

	sig := <-signalChan
	l.Println("Recieved terminate, graceful shutdown with sigtype:", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(tc)
}
