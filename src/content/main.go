package main

import (
	"fmt"
	"net"
	"os"
	"saltgram/content/gdrive"
	"saltgram/content/grpc/servers"
	"saltgram/content/saga"
	"saltgram/log"
	"saltgram/pki"
	"saltgram/protos/content/prcontent"
	"saltgram/protos/notifications/prnotifications"

	"saltgram/content/data"
	"saltgram/internal"

	"google.golang.org/grpc/reflection"
)

func main() {
	l := log.NewLogger("saltgram-content")
	l.L.Printf("Starting Content microservice on port: %s\n", os.Getenv("SALT_CONTENT_PORT"))
	pkiHandler := pki.Init()
	cert, err := pkiHandler.RegisterSaltgramService("saltgram-contents")
	if err != nil {
		l.L.Fatalf("failure while registering pki: %v\n", err)
	}
	s := internal.NewService(l.L)
	err = s.Init("saltgram-contents", cert.CertPEM, cert.PrivateKeyPEM, pkiHandler.RootCA.CertPEM)
	if err != nil {
		l.L.Fatalf("failure while initializing saltgram-contents: %v\n", err)
	}
	db := data.NewDBConn(l.L)
	db.ConnectToDb()
	db.MigradeData()

	g := gdrive.NewGDrive(l.L)

	rs := saga.NerRedisClient(l.L, db, g)
	go rs.Connection()
	nconn, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_NOTIF_ADDR", "localhost"), os.Getenv("SALT_NOTIF_PORT")))
	if err != nil {
		l.L.Fatalf("failure dialing notification: %v\n", err)
	}
	nc := prnotifications.NewNotificationsClient(nconn)

	gContentServer := servers.NewContent(l.L, db, g, nc)
	grpcServer := s.NewServer()
	prcontent.RegisterContentServer(grpcServer, gContentServer)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SALT_CONTENT_PORT")))
	if err != nil {
		l.L.Fatalf("failure while creating listener: %v\n", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		l.L.Fatalf("failure while serving: %v\n", err)
	}
	grpcServer.GracefulStop()
}
