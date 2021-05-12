package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"saltgram/auth/data"
	"saltgram/auth/grpc/servers"
	"saltgram/internal"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/users/prusers"

	"github.com/casbin/casbin/v2"
	"google.golang.org/grpc/reflection"
)

func main() {
	l := log.New(os.Stdout, "saltgram-auth", log.LstdFlags)
	l.Printf("Starting Auth microservice on port: %s\n", os.Getenv("SALT_AUTH_PORT"))
	s := internal.NewService(l)
	db := data.DBConn{}
	db.ConnectToDb()
	db.MigradeData()
	authEnforcer, err := casbin.NewEnforcer("./config/model.conf", "./config/policy.csv")
	if err != nil {
		l.Printf("[ERROR] creating auth enforcer: %v\n", err)
	}
	err = s.TLS.Init("../../certs/localhost.crt", "../../certs/localhost.key", "../../certs/RootCA.pem")
	if err != nil {
		l.Fatalf("[ERROR] configuring TLS: %v\n", err)
	}

	grpcServer := s.NewServer()
	usersConnection, err := s.GetConnection(fmt.Sprintf("localhost:%s", os.Getenv("SALT_USERS_PORT")))
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
}
