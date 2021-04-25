package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"saltgram/auth/data"
	"saltgram/auth/grpc/servers"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/users/prusers"

	"github.com/casbin/casbin/v2"
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
	l := log.New(os.Stdout, "saltgram-auth", log.LstdFlags)

	l.Printf("Starting Auth microservice on port: %s\n", os.Getenv("SALT_AUTH_PORT"))
	db := data.DBConn{}
	db.ConnectToDb()
	db.MigradeData()
	authEnforcer, err := casbin.NewEnforcer("./config/model.conf", "./config/policy.csv")
	if err != nil {
		l.Printf("[ERROR] creating auth enforcer: %v\n", err)
	}
	tlsConfig, err := getTLSConfig()
	if err != nil {
		l.Fatalf("[ERROR] configuring TLS: %v\n", err)
	}

	creds := credentials.NewTLS(tlsConfig)
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	usersConnection, err := grpc.Dial(fmt.Sprintf("localhost:%s", os.Getenv("SALT_USERS_PORT")), grpc.WithTransportCredentials(creds))
	if err != nil {
		l.Fatalf("[ERROR] dialing users connection: %v\n", err)
	}
	defer usersConnection.Close()
	usersClient := prusers.NewUsersClient(usersConnection)
	gAuthServer := servers.NewAuth(l, authEnforcer, &db, usersClient)
	prauth.RegisterAuthServer(grpcServer, gAuthServer)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SALT_AUTH_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] creating listener: %v\n", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		l.Fatalf("[ERROR] while serving: %v\n", err)
	}
	grpcServer.GracefulStop()
	// serverMux := mux.NewRouter()
	// authHandler := handlers.NewAuth(l)
	// jwtRouter := serverMux.PathPrefix("/jwt").Subrouter()
	// jwtRouter.HandleFunc("", authHandler.GetJWT(&db)).Methods(http.MethodPost)
	// jwtRouter.HandleFunc("", authHandler.Refresh(&db)).Methods(http.MethodGet)

	// refreshRouter := serverMux.PathPrefix("/refresh").Subrouter()
	// refreshRouter.HandleFunc("", authHandler.AddRefreshToken(&db)).Methods(http.MethodPost)

	// permRouter := serverMux.PathPrefix("/perm").Subrouter()
	// permRouter.HandleFunc("", authHandler.CheckPermissions(authEnforcer)).Methods(http.MethodGet)

	// loginHandler := handlers.NewLogin(l)
	// loginRouter := serverMux.PathPrefix("/login").Subrouter()
	// loginRouter.HandleFunc("", loginHandler.Login).Methods(http.MethodPost)
	// loginRouter.Use(loginHandler.MiddlewareValidateToken)

	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{fmt.Sprintf("https://localhost:%s", os.Getenv("SALT_API_PORT"))},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	// 	AllowCredentials: true,
	// 	Debug:            true,
	// })

	// server := http.Server{
	// 	Addr:         fmt.Sprintf(":%s", os.Getenv("SALT_AUTH_PORT")),
	// 	ReadTimeout:  1 * time.Second,
	// 	WriteTimeout: 1 * time.Second,
	// 	IdleTimeout:  120 * time.Second,
	// 	Handler:      c.Handler(serverMux),
	// 	TLSConfig:    tlsConfig,
	// }

	// go func() {
	// 	err := server.ListenAndServeTLS("", "")
	// 	if err != nil {
	// 		l.Fatalf("[ERROR] serving TLS: %v\n", err)
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
