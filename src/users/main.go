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
	"saltgram/users/data"
	"saltgram/users/handlers"
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

	passRouter := serverMux.PathPrefix("/changepw").Subrouter()
	passRouter.HandleFunc("", usersHandler.ChangePassword(&db)).Methods(http.MethodPost)
	passRouter.Use(usersHandler.MiddlewareValidateChangeRequest)

	emailRouter := serverMux.PathPrefix("/verifyemail").Subrouter()
	emailRouter.HandleFunc("", usersHandler.VerifyEmail(&db)).Methods(http.MethodPost)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://localhost:8080"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowCredentials: true,
		Debug:            true,
	})

	tlsConfig, err := getTLSConfig()
	if err != nil {
		l.Fatalf("[ERROR] configuring TLS: %v\n", err)
	}

	server := http.Server{
		Addr: fmt.Sprintf(":%s", os.Getenv("SALT_USERS_PORT")),
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
