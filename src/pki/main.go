package main

import (
	"log"
	"os"
	"saltgram/pki/data"
)

func main() {
	l := log.New(os.Stdout, "saltgram-pki", log.LstdFlags)
	db := data.NewDBConn(l)
	db.ConnectToDB()
	db.MigrateData()
	err := data.Init(db)
	if err != nil {
		l.Printf("Error initializing keystore: %v\n", err)
	}

	// name := pkix.Name{
	// 	CommonName: "Saltgram",
	// 	Organization: []string{"Saltgram"},
	// 	Country: []string{"RS"},
	// 	Province: []string{""},
	// 	Locality: []string{"Novi Sad"},
	// 	StreetAddress: []string{"Balzakova 69"},
	// 	PostalCode: []string{"21000"},
	// }
}