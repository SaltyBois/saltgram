package main

import (
	"log"
	"net/http"
	"os"
	"saltgram/users/data"
	"saltgram/users/handlers"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "saltgram-users", log.LstdFlags)

	db := data.DBConn{}
	db.ConnectToDb()
	db.MigradeData()

	usersHandler := handlers.NewUsers(l)

	serverMux := mux.NewRouter()
	getRouter := serverMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", usersHandler.GetByJWS(&db))

	postRouter := serverMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/users", usersHandler.Register(&db))
	postRouter.Use(usersHandler.MiddlewareValidateUser)

	putRouter := serverMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/users/{id:[0-9]+}", usersHandler.Update(&db))
	putRouter.Use(usersHandler.MiddlewareValidateUser)
}
