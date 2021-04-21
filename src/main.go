package main

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"saltgram/data"
	"saltgram/handlers"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	l := log.New(os.Stdout, "saltgram", log.LstdFlags)

	data.Seed()

	serverMux := mux.NewRouter()

	usersHandler := handlers.NewUsers(l)
	getRouter := serverMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", usersHandler.GetByJWS)
	// getRouter.HandleFunc("/users/{jws}", usersHandler.GetByJWS)
	postRouter := serverMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/users", usersHandler.Register)
	postRouter.Use(usersHandler.MiddlewareValidateUser)
	putRouter := serverMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/users/{id:[0-9]+}", usersHandler.Update)
	putRouter.Use(usersHandler.MiddlewareValidateUser)

	loginHandler := handlers.NewLogin(l)
	loginRouter := serverMux.PathPrefix("/login").Subrouter()
	loginRouter.HandleFunc("", loginHandler.Login).Methods(http.MethodPost)
	loginRouter.Use(loginHandler.MiddlewareValidateToken)

	authHandler := handlers.NewAuth(l)
	authRouter := serverMux.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/jwt", authHandler.GetJWT).Methods(http.MethodPost)
	authRouter.HandleFunc("/refresh", authHandler.Refresh).Methods(http.MethodGet)
	authRouter.HandleFunc("", authHandler.Logout).Methods(http.MethodDelete)
	// TODO(Jovan): Midleware?

	emailHandler := handlers.NewEmail(l)
	emailRouter := serverMux.PathPrefix("/email").Subrouter()
	emailRouter.HandleFunc("/activate/{token:[A-Za-z0-9]+}", emailHandler.Activate).Methods(http.MethodPut)
	emailRouter.HandleFunc("/activate", emailHandler.GetAll).Methods(http.MethodGet)
	emailRouter.HandleFunc("/change/{token:[A-Za-z0-9]+}", emailHandler.ConfirmReset).Methods(http.MethodPut)
	emailRouter.HandleFunc("/change", emailHandler.ChangePassword).Methods(http.MethodPost)
	emailRouter.HandleFunc("/forgot", emailHandler.RequestReset).Methods(http.MethodPost)

	// NOTE(Jovan): CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
		Debug:            true,
	})

	tlsConfig, err := getTLSConfig()
	if err != nil {
		l.Fatalf("[ERROR] loading TLS config: %v\n", err)
	}

	server := &http.Server{
		Addr:         os.Getenv("PORT_SALT"),
		Handler:      c.Handler(serverMux),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		TLSConfig:    tlsConfig,
	}

	go func() {
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatal(err)
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
