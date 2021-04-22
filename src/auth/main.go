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
	"saltgram/auth/data"
	"saltgram/auth/handlers"
	"time"

	"github.com/casbin/casbin/v2"
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
	l := log.New(os.Stdout, "saltgram-auth", log.LstdFlags)

	db := data.DBConn{}
	db.ConnectToDb()

	authEnforcer, err := casbin.NewEnforcer("./config/model.conf", "./config/policy.csv")
	if err != nil {
		l.Printf("[ERROR] creating auth enforcer: %v\n", err)
	}

	serverMux := mux.NewRouter()
	authHandler := handlers.NewAuth(l)
	jwtRouter := serverMux.PathPrefix("/jwt").Subrouter()
	jwtRouter.HandleFunc("", authHandler.GetJWT(&db)).Methods(http.MethodPost)
	jwtRouter.HandleFunc("/refresh", authHandler.Refresh(&db)).Methods(http.MethodGet)

	refreshRouter := serverMux.PathPrefix("/refresh").Subrouter()
	refreshRouter.HandleFunc("", authHandler.AddRefreshToken(&db)).Methods(http.MethodPost)

	permRouter := serverMux.PathPrefix("/perm").Subrouter()
	permRouter.HandleFunc("", authHandler.CheckPermissions(authEnforcer)).Methods(http.MethodGet)

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
		Addr: fmt.Sprintf(":%s", os.Getenv("SALT_AUTH_PORT")),
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout: 120 * time.Second,
		Handler: c.Handler(serverMux),
		TLSConfig: tlsConfig,
	}

	go func() {
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			l.Fatalf("[ERROR] serving TLS: %v\n", err)
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
