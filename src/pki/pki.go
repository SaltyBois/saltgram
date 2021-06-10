package pki

import (
	"crypto/x509/pkix"
	"log"
	"os"
	"saltgram/pki/data"
)

type PKI struct {
	db *data.DBConn
	RootCA *data.Certificate
}

func (p *PKI) RegisterSaltgramService(commonName string) (*data.Certificate, error) {
	return data.RegisterService(p.db, pkix.Name{
		CommonName: commonName,
		Organization: []string{"Saltgram"},
		Country: []string{"RS"},
		Province: []string{""},
		Locality: []string{"Liman"},
		StreetAddress: []string{"Balzakova 69"},
		PostalCode: []string{"21000"},
	})
}

func Init() *PKI {
	l := log.New(os.Stdout, "saltgram-pki", log.LstdFlags)
	db := data.NewDBConn(l)
	db.ConnectToDB()
	db.MigrateData()
	cert, err := data.Init(db)
	if err != nil {
		l.Fatalf("Error initializing keystore: %v\n", err)
	}
	return &PKI{db: db, RootCA: cert}
}