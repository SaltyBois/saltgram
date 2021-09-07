package main

import (
	"fmt"
	"net"
	"os"
	"saltgram/internal"
	"saltgram/log"
	"saltgram/notification/data"
	"saltgram/notification/grpc/servers"
	"saltgram/notification/pusher"
	"saltgram/pki"
	"saltgram/protos/notifications/prnotifications"
	"saltgram/protos/users/prusers"

	"google.golang.org/grpc/reflection"
)

func main() {
	l := log.NewLogger("saltgram-notification")
	l.L.Printf("Starting Notification microservice on port: %s\n", os.Getenv("SALT_NOTIF_PORT"))
	pkiHandler := pki.Init()
	cert, err := pkiHandler.RegisterSaltgramService("saltgram-notification")
	if err != nil {
		l.L.Fatalf("failure while registering pki: %v\n", err)
	}
	s := internal.NewService(l.L)
	err = s.Init("saltgram-notification", cert.CertPEM, cert.PrivateKeyPEM, pkiHandler.RootCA.CertPEM)
	if err != nil {
		l.L.Fatalf("failure while initializing saltgram-notification: %v\n", err)
	}
	db := data.NewDBConn(l.L)
	db.ConnectToDb()
	db.MigradeData()

	usersConnection, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_USERS_ADDR", "localhost"), os.Getenv("SALT_USERS_PORT")))
	if err != nil {
		l.L.Fatalf("dialing users connection: %v\n", err)
	}
	defer usersConnection.Close()
	usersClient := prusers.NewUsersClient(usersConnection)

	p := pusher.NewPusher(l.L)

	gNotificationServer := servers.NewNotification(l.L, db, usersClient, p)
	grpcServer := s.NewServer()
	prnotifications.RegisterNotificationsServer(grpcServer, gNotificationServer)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SALT_NOTIF_PORT")))
	if err != nil {
		l.L.Fatalf("failure while creating listener: %v\n", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		l.L.Fatalf("failure while serving: %v\n", err)
	}
	grpcServer.GracefulStop()

}
