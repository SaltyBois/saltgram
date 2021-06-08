package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"saltgram/content/grpc/servers"
	"saltgram/protos/content/prcontent"

	"saltgram/content/data"
	"saltgram/internal"

	"google.golang.org/grpc/reflection"
)

func main() {
	l := log.New(os.Stdout, "saltgram-contents", log.LstdFlags)
	l.Printf("Starting Content microservice on port: %s\n", os.Getenv("SALT_CONTENT_PORT"))
	s := internal.NewService(l)

	db := data.DBConn{}
	db.ConnectToDb()
	db.MigradeData()
	err := s.TLS.Init("../../certs/localhost.crt", "../../certs/localhost.key", "../../certs/RootCA.pem")
	if err != nil {
		l.Fatalf("[ERROR] configuring TLS: %v\n", err)
	}

	gContentServer := servers.NewContent(l, &db)
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
