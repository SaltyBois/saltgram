package pki

import (
	"agent/internal"
	"agent/pki/data"
	"crypto/x509/pkix"
)

type PKI struct {
	db     *data.DBConn
	RootCA *data.Certificate
}

func (p *PKI) RegisterSaltgramService(commonName string) (*data.Certificate, error) {
	return data.RegisterService(p.db, pkix.Name{
		CommonName:    commonName,
		Organization:  []string{"Saltgram"},
		Country:       []string{"RS"},
		Province:      []string{""},
		Locality:      []string{"Liman"},
		StreetAddress: []string{"Balzakova 69"},
		PostalCode:    []string{"21000"},
	})
}

func Init() *PKI {
	l := internal.NewLogger("saltgram-pki")
	db := data.NewDBConn(l.L)
	db.ConnectToDB()
	db.MigrateData()
	cert, err := data.Init(db)
	if err != nil {
		l.L.Fatalf("failure initializing keystore: %v\n", err)
	}
	return &PKI{db: db, RootCA: cert}
}
