package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"saltgram/internal"
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
	s := internal.NewService(l)

	db := data.NewDBConn(l)
	db.ConnectToDb()
	db.MigradeData()
	err := s.TLS.Init("../../certs/localhost.crt", "../../certs/localhost.key", "../../certs/RootCA.pem")
	if err != nil {
		l.Fatalf("[ERROR] configuring TLS: %v\n", err)
	}

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
