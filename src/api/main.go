package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"saltgram/api/handlers"
	"saltgram/internal"
	"saltgram/pki"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/content/prcontent"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"
	"time"

	"github.com/rs/cors"
)

func main() {
	l := log.New(os.Stdout, "saltgram-api-gateway", log.LstdFlags)
	l.Printf("Starting API Gateway on port: %s\n", os.Getenv("SALT_API_PORT"))
	pkiHandler := pki.Init()
	cert, err := pkiHandler.RegisterSaltgramService("api-gateway")
	if err != nil {
		l.Fatalf("[ERROR] registering service for pki: %v\n", err)
	}
	s := internal.NewService(l)
	
	err = s.Init("saltgram-api-gateway", cert.CertPEM, cert.PrivateKeyPEM, pkiHandler.RootCA.CertPEM)
	if err != nil {
		l.Fatalf("[ERROR] failed to init api service: %v\n", err)
	}

	authConnection, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_AUTH_ADDR", "localhost"), os.Getenv("SALT_AUTH_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] dialing auth connection: %v\n", err)
	}
	defer authConnection.Close()
	authClient := prauth.NewAuthClient(authConnection)
	authHandler := handlers.NewAuth(l, authClient)
	authRouter := s.S.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/refresh", authHandler.Refresh) //.Methods(http.MethodGet)
	authRouter.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	authRouter.HandleFunc("/jwt", authHandler.GetJWT).Methods(http.MethodPost)
	authRouter.HandleFunc("", authHandler.CheckPermissions).Methods(http.MethodPut)

	usersConnection, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_USERS_ADDR", "localhost"), os.Getenv("SALT_USERS_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] dialing users connection: %v\n", err)
	}
	defer usersConnection.Close()
	usersClient := prusers.NewUsersClient(usersConnection)
	usersHandler := handlers.NewUsers(l, usersClient)
	usersRouter := s.S.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/register", usersHandler.Register).Methods(http.MethodPost)
	usersRouter.HandleFunc("", usersHandler.GetByJWS).Methods(http.MethodGet)
	usersRouter.HandleFunc("/resetpass", usersHandler.ResetPassword).Methods(http.MethodPost)
	usersRouter.HandleFunc("/changepass", usersHandler.ChangePassword).Methods(http.MethodPost)

	emailConnection, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_EMAIL_ADDR", "localhost"), os.Getenv("SALT_EMAIL_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] dialing email connection")
	}

	emailClient := premail.NewEmailClient(emailConnection)
	emailHandler := handlers.NewEmail(l, emailClient, usersClient)
	emailRouter := s.S.PathPrefix("/email").Subrouter()
	emailRouter.HandleFunc("/activate/{token}", emailHandler.Activate).Methods(http.MethodPut)
	emailRouter.HandleFunc("/forgot", emailHandler.ForgotPassword).Methods(http.MethodPost)
	emailRouter.HandleFunc("/reset/{token}", emailHandler.ConfirmReset).Methods(http.MethodPut)

	contentConnection, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_CONTENT_ADDR", "localhost"), os.Getenv("SALT_CONTENT_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] dialing content connection: %v\n", err)
	}
	defer contentConnection.Close()
	contentClient := prcontent.NewContentClient(contentConnection)
	contentHandler := handlers.NewContent(l, contentClient, usersClient)
	contentRouter := s.S.PathPrefix("/content").Subrouter()
	contentRouter.HandleFunc("/user", contentHandler.GetSharedMedia).Methods(http.MethodGet)
	contentRouter.HandleFunc("/sharedmedia", contentHandler.AddSharedMedia).Methods(http.MethodPost)
	contentRouter.HandleFunc("/user/{id}", contentHandler.GetSharedMediaByUser).Methods(http.MethodGet)

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
		Handler:      c.Handler(s.S),
		TLSConfig:    s.TLS.TC,
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
