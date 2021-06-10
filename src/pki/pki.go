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
	// s := internal.NewService(l)
	// err = s.Init(cert.Cert.Subject.CommonName, cert.CertPEM, cert.PrivateKeyPEM, cert.CertPEM)
	// if err != nil {
	// 	l.Fatalf("[ERROR] initializing keystore: %v\n", err)
	// }
	// grpcServer := s.NewServer()
	// pkiServer := servers.NewPKI(l, db)
	// prpki.RegisterPKIServer(grpcServer, pkiServer)
	// listener, err := net.Listen("tcp", fmt.Sprintf(":%s", internal.GetEnvOrDefault("SALT_PKI_PORT", "8086")))
	// if err != nil {
	// 	l.Fatalf("[ERROR] failed to start PKI server: %v\n", err)
	// }
	// err = grpcServer.Serve(listener)
	// if err != nil {
	// 	l.Fatalf("[ERROR] failed to serve PKI server: %v\n", err)
	// }
	// grpcServer.GracefulStop()

}