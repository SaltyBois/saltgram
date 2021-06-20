package main

import (
	"fmt"
	"net"
	"os"
	"saltgram/internal"
	"saltgram/log"
	"saltgram/pki"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/content/prcontent"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"
	"saltgram/users/data"
	"saltgram/users/grpc/servers"

	"google.golang.org/grpc/reflection"
)

func main() {
	l := log.NewLogger("saltgram-users")
	l.L.Infof("Starting Users microservice on port: %s\n", os.Getenv("SALT_USERS_PORT"))
	pkiHandler := pki.Init()
	cert, err := pkiHandler.RegisterSaltgramService("saltgram-users")
	if err != nil {
		l.L.Fatalf("failure while registering to PKI: %v\n", err)
	}
	s := internal.NewService(l.L)
	err = s.Init("saltgram-users", cert.CertPEM, cert.PrivateKeyPEM, pkiHandler.RootCA.CertPEM)
	if err != nil {
		l.L.Fatalf("failure initializing users service: %v\n", err)
	}
	db := data.NewDBConn(l.L)
	db.ConnectToDb()
	db.MigradeData()

	aconn, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_AUTH_ADDR", "localhost"), os.Getenv("SALT_AUTH_PORT")))
	if err != nil {
		l.L.Fatalf("failure dialing auth: %v\n", err)
	}
	ac := prauth.NewAuthClient(aconn)
	econn, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_EMAIL_ADDR", "localhost"), os.Getenv("SALT_EMAIL_PORT")))
	if err != nil {
		l.L.Fatalf("failure dialing email: %v\n", err)
	}
	ec := premail.NewEmailClient(econn)

	cconn, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_CONTENT_ADDR", "localhost"), os.Getenv("SALT_CONTENT_PORT")))
	if err != nil {
		l.L.Fatalf("failure dialing content: %v", err)
	}
	cc := prcontent.NewContentClient(cconn)

	gUsersServer := servers.NewUsers(l.L, db, ac, ec, cc)
	grpcServer := s.NewServer()
	prusers.RegisterUsersServer(grpcServer, gUsersServer)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SALT_USERS_PORT")))
	if err != nil {
		l.L.Fatalf("failure creating listener: %v\n", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		l.L.Fatalf("failure while serving: %v\n", err)
	}
	grpcServer.GracefulStop()
}
