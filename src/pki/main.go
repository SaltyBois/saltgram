package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"saltgram/internal"
	"saltgram/pki/data"
	"saltgram/pki/grpc/servers"
	"saltgram/protos/pki/prpki"
)

func main() {
	l := log.New(os.Stdout, "saltgram-pki", log.LstdFlags)
	db := data.NewDBConn(l)
	db.ConnectToDB()
	db.MigrateData()
	cert, err := data.Init(db)
	if err != nil {
		l.Fatalf("Error initializing keystore: %v\n", err)
	}
	s := internal.NewService(l)
	err = s.Init(cert.Cert.Subject.CommonName, cert.CertPEM, cert.PrivateKeyPEM, cert.CertPEM)
	if err != nil {
		l.Fatalf("[ERROR] initializing keystore: %v\n", err)
	}
	grpcServer := s.NewServer()
	pkiServer := servers.NewPKI(l, db)
	prpki.RegisterPKIServer(grpcServer, pkiServer)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", internal.GetEnvOrDefault("SALT_PKI_PORT", "8086")))
	if err != nil {
		l.Fatalf("[ERROR] failed to start PKI server: %v\n", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		l.Fatalf("[ERROR] failed to serve PKI server: %v\n", err)
	}
	grpcServer.GracefulStop()

}