package main

import (
	"fmt"
	"net"
	"os"
	"saltgram/email/data"
	"saltgram/email/grpc/servers"
	"saltgram/internal"
	"saltgram/log"
	"saltgram/pki"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"

	"google.golang.org/grpc/reflection"
)

func main() {
	l := log.NewLogger("saltgram-email")
	l.L.Infof("Starting Email microservice on port: %s\n", os.Getenv("SALT_EMAIL_PORT"))
	pkiHandler := pki.Init()
	cert, err := pkiHandler.RegisterSaltgramService("saltgram-email")
	if err != nil {
		l.L.Fatalf("failure while registering pki: %v\n", err)
	}
	s := internal.NewService(l.L)
	err = s.Init("saltgram-email", cert.CertPEM, cert.PrivateKeyPEM, pkiHandler.RootCA.CertPEM)
	if err != nil {
		l.L.Fatalf("failed to initialize email service: %v\n", err)
	}
	db := data.NewDBConn(l.L)
	db.ConnectToDb()
	db.MigradeData()

	uconn, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_USERS_ADDR", "localhost"), os.Getenv("SALT_USERS_PORT")))
	if err != nil {
		l.L.Fatalf("failure dialing users connection: %v\n", err)
	}
	defer uconn.Close()

	uc := prusers.NewUsersClient(uconn)
	emailServer := servers.NewEmail(l.L, db, uc)
	grpcServer := s.NewServer()
	premail.RegisterEmailServer(grpcServer, emailServer)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SALT_EMAIL_PORT")))
	if err != nil {
		l.L.Fatalf("failure creating listener: %v\n", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		l.L.Fatalf("failure while serving: %v\n", err)
	}
	grpcServer.GracefulStop()
}
