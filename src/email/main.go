package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"saltgram/email/data"
	"saltgram/email/grpc/servers"
	"saltgram/internal"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"

	"google.golang.org/grpc/reflection"
)

func main() {
	l := log.New(os.Stdout, "saltgram-email", log.LstdFlags)
	l.Printf("Starting Email microservice on port: %s\n", os.Getenv("SALT_EMAIL_PORT"))
	s := internal.NewService(l)
	db := data.DBConn{}
	db.ConnectToDb()
	db.MigradeData()
	err := s.TLS.Init("../../certs/localhost.crt", "../../certs/localhost.key", "../../certs/RootCA.pem")
	if err != nil {
		l.Fatalf("[ERROR] configuring tls: %v\n", err)
	}

	uconn, err := s.GetConnection(fmt.Sprintf("localhost:%s", os.Getenv("SALT_USERS_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] dialing users connection: %v\n", err)
	}
	defer uconn.Close()

	uc := prusers.NewUsersClient(uconn)
	emailServer := servers.NewEmail(l, &db, uc)
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
