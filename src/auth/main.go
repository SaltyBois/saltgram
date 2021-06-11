package main

import (
	"fmt"
	"net"
	"os"
	"saltgram/auth/data"
	"saltgram/auth/grpc/servers"
	"saltgram/internal"
	"saltgram/log"
	"saltgram/pki"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/users/prusers"

	"github.com/casbin/casbin/v2"
	"google.golang.org/grpc/reflection"
)

func main() {
	l := log.NewLogger("saltgram-auth")
	l.L.Printf("Starting Auth microservice on port: %s\n", os.Getenv("SALT_AUTH_PORT"))
	pkiHandler := pki.Init()
	cert, err := pkiHandler.RegisterSaltgramService("saltgram-auth")
	if err != nil {
		l.L.Fatalf("failed to register to PKI")
	}
	s := internal.NewService(l.L)
	err = s.Init("saltgram-auth", cert.CertPEM, cert.PrivateKeyPEM, pkiHandler.RootCA.CertPEM)
	if err != nil {
		l.L.Fatalf("failed to init auth service: %v\n", err)
	}
	db := data.NewDBConn(l.L)
	db.ConnectToDb()
	db.MigradeData()
	authEnforcer, err := casbin.NewEnforcer("./config/model.conf", "./config/policy.csv")
	if err != nil {
		l.L.Errorf("failure creating auth enforcer: %v\n", err)
	}

	grpcServer := s.NewServer()
	usersConnection, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_USERS_ADDR", "localhost"), os.Getenv("SALT_USERS_PORT")))
	if err != nil {
		l.L.Fatalf("failure dialing users connection: %v\n", err)
	}
	defer usersConnection.Close()

	usersClient := prusers.NewUsersClient(usersConnection)
	gAuthServer := servers.NewAuth(l.L, authEnforcer, db, usersClient)
	prauth.RegisterAuthServer(grpcServer, gAuthServer)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SALT_AUTH_PORT")))
	if err != nil {
		l.L.Fatalf("failure creating listener: %v\n", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		l.L.Fatalf("failure while serving: %v\n", err)
	}
	grpcServer.GracefulStop()
}
