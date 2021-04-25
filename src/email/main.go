package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"saltgram/email/data"
	"saltgram/email/grpc/servers"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"

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
	l := log.New(os.Stdout, "saltgram-email", log.LstdFlags)
	l.Printf("Starting Email microservice on port: %s\n", os.Getenv("SALT_EMAIL_PORT"))
	db := data.DBConn{}
	db.ConnectToDb()
	db.MigradeData()
	tlsConfig, err := getTLSConfig()
	if err != nil {
		l.Fatalf("[ERROR] configuring tls: %v\n", err)
	}
	creds := credentials.NewTLS(tlsConfig)

	uconn, err := grpc.Dial(fmt.Sprintf("localhost:%s", os.Getenv("SALT_USERS_PORT")), grpc.WithTransportCredentials(creds))
	if err != nil {
		l.Fatalf("[ERROR] dialing users connection: %v\n", err)
	}
	defer uconn.Close()
	uc := prusers.NewUsersClient(uconn)
	emailServer := servers.NewEmail(l, &db, uc)
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	premail.RegisterEmailServer(grpcServer, emailServer)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SALT_EMAIL_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] creating listener: %v\n", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		l.Fatalf("[ERROR] while serving: %v\n", err)
	}
	grpcServer.GracefulStop()

	// serverMux := mux.NewRouter()

	// activationRouter := serverMux.PathPrefix("/activate").Subrouter()
	// activationRouter.HandleFunc("", emailHandler.SendActivation).Methods(http.MethodPost)
	// activationRouter.HandleFunc("/{token:[A-Za-z0-9]+}", emailHandler.Activate).Methods(http.MethodPut)

	// changeRouter := serverMux.PathPrefix("/change").Subrouter()
	// changeRouter.HandleFunc("/{token:[A-Za-z0-9]+}", emailHandler.ConfirmReset).Methods(http.MethodPut)
	// changeRouter.HandleFunc("", emailHandler.ChangePassword).Methods(http.MethodPost)
	// changeRouter.HandleFunc("/forgot", emailHandler.RequestReset).Methods(http.MethodPost)

	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{fmt.Sprintf("https://localhost:%s", os.Getenv("SALT_API_PORT"))},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
	// 	AllowCredentials: true,
	// 	Debug:            true,
	// })

	// server := http.Server{
	// 	Addr:         fmt.Sprintf(":%s", os.Getenv("SALT_EMAIL_PORT")),
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
