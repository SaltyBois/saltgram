package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"saltgram/email/data"
	"saltgram/email/grpc/servers"
	"saltgram/internal"
	"saltgram/pki"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"

	"google.golang.org/grpc/reflection"
)

func main() {
	l := log.New(os.Stdout, "saltgram-email", log.LstdFlags)
	l.Printf("Starting Email microservice on port: %s\n", os.Getenv("SALT_EMAIL_PORT"))
	pkiHandler := pki.Init()
	cert, err := pkiHandler.RegisterSaltgramService("saltgram-email")
	if err != nil {
		l.Fatalf("[ERROR] while registering pki: %v\n", err)
	}
	s := internal.NewService(l)
	err = s.Init("saltgram-email", cert.CertPEM, cert.PrivateKeyPEM, pkiHandler.RootCA.CertPEM)
	db := data.NewDBConn(l)
	db.ConnectToDb()
	db.MigradeData()

	uconn, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_USERS_ADDR", "localhost"), os.Getenv("SALT_USERS_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] dialing users connection: %v\n", err)
	}
	defer uconn.Close()

	uc := prusers.NewUsersClient(uconn)
	emailServer := servers.NewEmail(l, db, uc)
	grpcServer := s.NewServer()
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
}
