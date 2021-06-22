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
	"saltgram/protos/content/prcontent"

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

	contentConnection, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_CONTENT_ADDR", "localhost"), os.Getenv("SALT_CONTENT_PORT")))
	if err != nil {
		l.L.Fatalf("dialing content connection: %v\n", err)
	}
	defer contentConnection.Close()
	contentClient := prcontent.NewContentClient(contentConnection)

	gAdminServer := servers.NewAdmin(l.L, db, contentClient)
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
