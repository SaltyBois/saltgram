package main

import (
	"fmt"
	"net"
	"os"
	"saltgram/admin/grpc/servers"
	"saltgram/log"
	"saltgram/pki"

	"saltgram/admin/data"
	"saltgram/internal"
	"saltgram/protos/admin/pradmin"

	"google.golang.org/grpc/reflection"
)

func main() {
	l := log.NewLogger("saltgram-admin")
	l.L.Printf("Starting Admin microservice on port: %s\n", os.Getenv("SALT_ADMIN_PORT"))
	pkiHandler := pki.Init()
	cert, err := pkiHandler.RegisterSaltgramService("saltgram-admin")
	if err != nil {
		l.L.Fatalf("failure while registering pki: %v\n", err)
	}
	s := internal.NewService(l.L)
	err = s.Init("saltgram-admin", cert.CertPEM, cert.PrivateKeyPEM, pkiHandler.RootCA.CertPEM)
	if err != nil {
		l.L.Fatalf("failure while initializing saltgram-admin: %v\n", err)
	}
	db := data.NewDBConn(l.L)
	db.ConnectToDb()
	db.MigradeData()

	gAdminServer := servers.NewAdmin(l.L, db)
	grpcServer := s.NewServer()
	pradmin.RegisterAdminServer(grpcServer, gAdminServer)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SALT_ADMIN_PORT")))
	if err != nil {
		l.L.Fatalf("failure while creating listener: %v\n", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		l.L.Fatalf("failure while serving: %v\n", err)
	}
	grpcServer.GracefulStop()
}
