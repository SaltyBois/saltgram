package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"
	"saltgram/users/data"
	"saltgram/users/grpc/servers"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
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

func main() {
	l := log.New(os.Stdout, "saltgram-users", log.LstdFlags)
	l.Printf("Starting Users microservice on port: %s\n", os.Getenv("SALT_USERS_PORT"))

	db := data.DBConn{}
	db.ConnectToDb()
	db.MigradeData()
	tlsConfig, err := getTLSConfig()
	if err != nil {
		l.Fatalf("[ERROR] configuring TLS: %v\n", err)
	}

	creds := credentials.NewTLS(tlsConfig)
	aconn, err := grpc.Dial(fmt.Sprintf("localhost:%s", os.Getenv("SALT_AUTH_PORT")), grpc.WithTransportCredentials(creds))
	if err != nil {
		l.Fatalf("[ERROR] dialing auth: %v\n", err)
	}
	ac := prauth.NewAuthClient(aconn)
	econn, err := grpc.Dial(fmt.Sprintf("localhost:%s", os.Getenv("SALT_EMAIL_PORT")), grpc.WithTransportCredentials(creds))
	if err != nil {
		l.Fatalf("[ERROR] dialing email: %v\n", err)
	}
	ec := premail.NewEmailClient(econn)
	gUsersServer := servers.NewUsers(l, &db, ac, ec)
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	prusers.RegisterUsersServer(grpcServer, gUsersServer)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SALT_USERS_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] creating listener: %v\n", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		l.Fatalf("[ERROR] while serving: %v\n", err)
	}
	grpcServer.GracefulStop()

	// usersHandler := handlers.NewUsers(l)

	// serverMux := mux.NewRouter()
	// getRouter := serverMux.Methods(http.MethodGet).Subrouter()
	// getRouter.HandleFunc("/users", usersHandler.GetByJWS(&db))

	// postRouter := serverMux.Methods(http.MethodPost).Subrouter()
	// postRouter.HandleFunc("/users", usersHandler.Register(&db))
	// postRouter.Use(usersHandler.MiddlewareValidateUser)

	// putRouter := serverMux.Methods(http.MethodPut).Subrouter()
	// putRouter.HandleFunc("/users/{id:[0-9]+}", usersHandler.Update(&db))
	// putRouter.Use(usersHandler.MiddlewareValidateUser)

	// passRouter := serverMux.PathPrefix("/password").Subrouter()
	// passRouter.HandleFunc("", usersHandler.ChangePassword(&db)).Methods(http.MethodPost)
	// passRouter.HandleFunc("", usersHandler.IsPasswordValid(&db)).Methods(http.MethodPut)
	// passRouter.Use(usersHandler.MiddlewareValidateChangeRequest)

	// emailRouter := serverMux.PathPrefix("/verifyemail").Subrouter()
	// emailRouter.HandleFunc("", usersHandler.VerifyEmail(&db)).Methods(http.MethodPost)
	// emailRouter.HandleFunc("/{un:[A-Za-z0-9]+}", usersHandler.IsEmailVerified(&db)).Methods(http.MethodGet)

	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{fmt.Sprintf("https://localhost:%s", os.Getenv("SALT_API_PORT"))},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	// 	AllowCredentials: true,
	// 	Debug:            true,
	// })

	// server := http.Server{
	// 	Addr:         fmt.Sprintf(":%s", os.Getenv("SALT_USERS_PORT")),
	// 	ReadTimeout:  1 * time.Second,
	// 	WriteTimeout: 1 * time.Second,
	// 	IdleTimeout:  120 * time.Second,
	// 	Handler:      c.Handler(serverMux),
	// 	TLSConfig:    tlsConfig,
	// }

	// go func() {
	// 	err := server.ListenAndServeTLS("", "")
	// 	if err != nil {
	// 		l.Fatalf("[ERROR] while serving: %v\n", err)
	// 	}
	// }()

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt, os.Kill)

	// sig := <-signalChan
	// l.Println("Recieved terminate, graceful shutdown with sigtype:", sig)

	// tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	// server.Shutdown(tc)
}
