package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"saltgram/internal"
	"saltgram/pki"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"
	"saltgram/users/data"
	"saltgram/users/grpc/servers"

	"google.golang.org/grpc/reflection"
)

func main() {
	l := log.New(os.Stdout, "saltgram-users", log.LstdFlags)
	l.Printf("Starting Users microservice on port: %s\n", os.Getenv("SALT_USERS_PORT"))
	pkiHandler := pki.Init()
	cert, err := pkiHandler.RegisterSaltgramService("saltgram-users")
	if err != nil {
		l.Fatalf("[ERROR] while registering to PKI: %v\n", err)
	}
	s := internal.NewService(l)
	err = s.Init("saltgram-users", cert.CertPEM, cert.PrivateKeyPEM, pkiHandler.RootCA.CertPEM)
	if err != nil {
		l.Fatalf("[ERROR] initializing users service: %v\n", err)
	}
	db := data.NewDBConn(l)
	db.ConnectToDb()
	db.MigradeData()

	aconn, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_AUTH_ADDR", "localhost"), os.Getenv("SALT_AUTH_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] dialing auth: %v\n", err)
	}
	ac := prauth.NewAuthClient(aconn)
	econn, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_EMAIL_ADDR", "localhost"), os.Getenv("SALT_EMAIL_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] dialing email: %v\n", err)
	}
	ec := premail.NewEmailClient(econn)
	gUsersServer := servers.NewUsers(l, db, ac, ec)
	grpcServer := s.NewServer()
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
}
