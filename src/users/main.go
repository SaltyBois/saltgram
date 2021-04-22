package main

import (
	"log"
	"net/http"
	"os"
	"saltgram/users/handlers"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "saltgram-users", log.LstdFlags)

	usersHandler := handlers.NewUsers(l)

	serverMux := mux.NewRouter()
	getRouter := serverMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", usersHandler.GetByJWS)

	postRouter := serverMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/users", usersHandler.Register)
	postRouter.Use(usersHandler.MiddlewareValidateUser)

	putRouter := serverMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/users/{id:[0-9]+}", usersHandler.Update)
	putRouter.Use(usersHandler.MiddlewareValidateUser)
}
