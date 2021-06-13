package main

import (
	"fmt"
	"net"
	"os"
	"saltgram/notification/grpc/servers"
	"saltgram/log"
	"saltgram/pki"
	"saltgram/protos/notification/prnotification"

	"saltgram/notification/data"
	"saltgram/internal"

	"google.golang.org/grpc/reflection"
)

func main() {
	l := log.NewLogger("saltgram-notification")
	l.L.Printf("Starting Notification microservice on port: %s\n", os.Getenv("SALT_NOTIFICATION_PORT"))
	pkiHandler := pki.Init()
	cert, err := pkiHandler.RegisterSaltgramService("saltgram-notifications")
	if err != nil {
		l.L.Fatalf("failure while registering pki: %v\n", err)
	}
	s := internal.NewService(l.L)
	err = s.Init("saltgram-notifications", cert.CertPEM, cert.PrivateKeyPEM, pkiHandler.RootCA.CertPEM)
	if err != nil {
		l.L.Fatalf("failure while initializing saltgram-notifications: %v\n", err)
	}
	db := data.NewDBConn(l.L)
	db.ConnectToDb()
	db.MigradeData()

	gNotificationServer := servers.NewNotification(l.L, db)
	grpcServer := s.NewServer()
	prnotification.RegisterNotificationServer(grpcServer, gNotificationServer)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SALT_NOTIFICATION_PORT")))
	if err != nil {
		l.L.Fatalf("failure while creating listener: %v\n", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		l.L.Fatalf("failure while serving: %v\n", err)
	}
	grpcServer.GracefulStop()
}