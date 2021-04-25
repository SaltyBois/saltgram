package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"saltgram/api/handlers"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func loadTLSCert(crtPath, keyPath string) (*tls.Certificate, error) {
	crt, err := ioutil.ReadFile(crtPath)
	if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

func getTLSConfig() (*tls.Config, error) {
	cert, err := loadTLSCert("../../certs/localhost.crt", "../../certs/localhost.key")
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile("../../certs/RootCA.pem")
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{*cert},
		ServerName:   "localhost",
		RootCAs:      caPool,
	}, nil
}

func getConnection(creds credentials.TransportCredentials, addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {
	l := log.New(os.Stdout, "saltgram-api-gateway", log.LstdFlags)
	l.Printf("Starting API Gateway on port: %s\n", os.Getenv("SALT_API_PORT"))

	serverMux := mux.NewRouter()

	tlsConfig, err := getTLSConfig()
	if err != nil {
		l.Fatalf("[ERROR] configuring tls: %v\n", err)
	}

	creds := credentials.NewTLS(tlsConfig)
	authConnection, err := getConnection(creds, fmt.Sprintf("localhost:%s", os.Getenv("SALT_AUTH_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] dialing auth connection: %v\n", err)
	}
	defer authConnection.Close()
	authClient := prauth.NewAuthClient(authConnection)
	authHandler := handlers.NewAuth(l, authClient)
	authRouter := serverMux.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)

	usersConnection, err := getConnection(creds, fmt.Sprintf("localhost:%s", os.Getenv("SALT_USERS_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] dialing users connection: %v\n", err)
	}
	defer usersConnection.Close()
	usersClient := prusers.NewUsersClient(usersConnection)
	usersHandler := handlers.NewUsers(l, usersClient)
	usersRouter := serverMux.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/register", usersHandler.Register).Methods(http.MethodPost)

	emailConnection, err := getConnection(creds, fmt.Sprintf("localhost:%s", os.Getenv("SALT_EMAIL_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] dialing email connection")
	}
	emailClient := premail.NewEmailClient(emailConnection)
	emailHandler := handlers.NewEmail(l, emailClient)
	emailRouter := serverMux.PathPrefix("/email").Subrouter()
	emailRouter.HandleFunc("/activate/{token}", emailHandler.Activate).Methods(http.MethodPut)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{fmt.Sprintf("https://localhost:%s", os.Getenv("SALT_WEB_PORT"))},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
		Debug:            true,
	})

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("SALT_API_PORT")),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      c.Handler(serverMux),
		TLSConfig:    tlsConfig,
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
