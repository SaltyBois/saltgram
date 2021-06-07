package main

import (
	"log"
	"os"

	"saltgram/internal"
	"saltgram/content/data"
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


}
