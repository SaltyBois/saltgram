package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"saltgram/content/grpc/servers"
	"saltgram/pki"
	"saltgram/protos/content/prcontent"

	"saltgram/content/data"
	"saltgram/internal"

	"google.golang.org/grpc/reflection"
)

func main() {
	l := log.New(os.Stdout, "saltgram-contents", log.LstdFlags)
	l.Printf("Starting Content microservice on port: %s\n", os.Getenv("SALT_CONTENT_PORT"))
	pkiHandler := pki.Init()
	cert, err := pkiHandler.RegisterSaltgramService("saltgram-contents")
	if err != nil {
		l.Fatalf("[ERROR] while registering pki: %v\n", err)
	}
	s := internal.NewService(l)
	err = s.Init("saltgram-contents", cert.CertPEM, cert.PrivateKeyPEM, pkiHandler.RootCA.CertPEM)
	if err != nil {
		l.Fatalf("[ERROR] while initializing saltgram-contents: %v\n", err)
	}
	db := data.NewDBConn(l)
	db.ConnectToDb()
	db.MigradeData()

	gContentServer := servers.NewContent(l, db)
	grpcServer := s.NewServer()
	prcontent.RegisterContentServer(grpcServer, gContentServer)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SALT_CONTENT_PORT")))
	if err != nil {
		l.Fatalf("[ERROR] creating listener: %v\n", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		l.Fatalf("[ERROR] while serving: %v\n", err)
	}
	grpcServer.GracefulStop()
}
