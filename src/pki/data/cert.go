package data

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	"saltgram/internal"
	"time"
)

var (
	ROOT_CN = internal.GetEnvOrDefault("PKI_ROOT_CN", "SaltgramRootCA")
)

type Certificate struct {
	Cert          *x509.Certificate   `json:"cert"`
	CertChain     []*x509.Certificate `json:"certChain"`
	PrivateKey    *rsa.PrivateKey     `json:"-"`
	Type          CertType            `json:"type"`
	CertPEM       []byte
	PrivateKeyPEM []byte
}

type ArchivedCert struct {
	SerialNumber string    `gorm:"primaryKey" json:"serialNumber"`
	ArchiveDate  time.Time `json:"archiveDate"`
}

type LookupDTO struct {
	SubjectName  string `json:"subjectName"`
	SerialNumber string `json:"serialNumber"`
}

var trustedDNS = []string{"localhost", "saltgram-auth", "saltgram-api-gateway", "saltgram-contents", "saltgram-email", "saltgram-users", "saltgram-webserver"}

func Init(db *DBConn) (*Certificate, error) {
	InitKeystore()
	cert, err := LoadCert(LookupDTO{
		SubjectName: ROOT_CN,
	})
	if err != nil {
		serialNumber := GetRandomSerial()
		rootTemplate := &x509.Certificate{
			SerialNumber: serialNumber,
			Subject: pkix.Name{
				SerialNumber:  serialNumber.String(),
				CommonName:    "SaltgramRootCA",
				Organization:  []string{"Saltgram"},
				Country:       []string{"RS"},
				Province:      []string{"Liman"},
				Locality:      []string{"Novi Sad"},
				StreetAddress: []string{"Balzakova 69"},
				PostalCode:    []string{"21000"},
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(10, 0, 0),
			IsCA:                  true,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			BasicConstraintsValid: true,
			DNSNames:              trustedDNS,
		}

		cert, err = GenRootCA(rootTemplate)
		if err != nil {
			return nil, err
		}

		err = cert.Save()
		if err != nil {
			return nil, err
		}
	}

	return cert, nil
}

func RegisterService(db *DBConn, subject pkix.Name) (*Certificate, error) {
	serialNumber := GetRandomSerial()
	template := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  false,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: false,
		DNSNames:              trustedDNS,
	}

	cert, err := GenCert(db, template, LookupDTO{
		SubjectName:  ROOT_CN,
		SerialNumber: "",
	})
	if err != nil {
		return nil, err
	}

	if err := cert.Save(); err != nil {
		return nil, err
	}

	cert, err = LoadCert(LookupDTO{SubjectName: subject.CommonName})
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func GenRootCA(rootTemplate *x509.Certificate) (*Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	rootCert, _, err := genCert(rootTemplate, rootTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}
	rootCert.Issuer = rootCert.Subject
	cert := Certificate{Cert: rootCert, PrivateKey: privateKey, Type: Root}
	return &cert, nil
}

func GenCert(db *DBConn, template *x509.Certificate, issuerDTO LookupDTO) (*Certificate, error) {
	issuerCert := Certificate{}
	if err := issuerCert.Load(issuerDTO); err != nil {
		return nil, err
	}

	if err := issuerCert.Verify(db); err != nil {
		return nil, err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	c, _, err := genCert(template, issuerCert.Cert, &privateKey.PublicKey, issuerCert.PrivateKey)
	if err != nil {
		return nil, err
	}
	certChain := issuerCert.CertChain
	certChain = append([]*x509.Certificate{issuerCert.Cert}, certChain...)

	cert := Certificate{
		Cert:       c,
		CertChain:  certChain,
		PrivateKey: privateKey,
		Type:       GetType(c),
	}
	return &cert, err
}

func GetType(c *x509.Certificate) CertType {
	if !c.IsCA {
		return EndEntity
	} else if c.SerialNumber.String() != c.Issuer.SerialNumber {
		return Intermediary
	} else {
		return Root
	}
}

func LoadCert(dto LookupDTO) (*Certificate, error) {
	cert := &Certificate{}
	err := cert.Load(dto)
	return cert, err
}

func (cert *Certificate) Load(dto LookupDTO) error {
	if err := cert.loadCertAndKey(dto.SubjectName); err != nil {
		return err
	}
	return nil
}

var ErrCertificateIsArchived = fmt.Errorf("certificate is archived")

func (cert *Certificate) Verify(db *DBConn) error {
	if cert.Type == Root && !IsArchived(db, []string{cert.Cert.Subject.SerialNumber}) {
		return nil
	}
	if IsArchived(db, func() []string {
		serialNumbers := []string{cert.Cert.SerialNumber.String()}
		for _, c := range cert.CertChain {
			serialNumbers = append(serialNumbers, c.Subject.SerialNumber)
		}
		return serialNumbers
	}()) {
		return ErrCertificateIsArchived
	}
	roots, inter := getRootInetrPool(cert)

	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: inter,
		CurrentTime:   time.Now(),
	}

	if _, err := cert.Cert.Verify(opts); err != nil {
		return err
	}
	return nil
}

func (cert *Certificate) Save() error {
	filename := cert.Cert.Subject.CommonName //+ cert.Cert.SerialNumber.String()
	err := WritePFX(cert.Cert, cert.CertChain, cert.PrivateKey, filename)
	if err != nil {
		return err
	}
	return nil
}

func ArchiveCert(db *DBConn, dto LookupDTO) error {
	cert := Certificate{}
	cert.Load(dto)
	archivedCert := ArchivedCert{SerialNumber: cert.Cert.SerialNumber.String(), ArchiveDate: time.Now()}
	return db.DB.Create(&archivedCert).Error
}

func IsArchived(db *DBConn, serialNumbers []string) bool {
	archive := ArchivedCert{}
	return db.DB.Find(&archive).Where("serial_number IN (?)", serialNumbers).RowsAffected > 0
}

func getRootInetrPool(cert *Certificate) (*x509.CertPool, *x509.CertPool) {
	root := x509.NewCertPool()
	inter := x509.NewCertPool()
	for _, c := range cert.CertChain {
		if GetType(c) == Root {
			root.AddCert(c)
		} else {
			inter.AddCert(c)
		}
	}
	return root, inter
}

func genCert(template, parent *x509.Certificate, subjectKey *rsa.PublicKey, issuerKey *rsa.PrivateKey) (*x509.Certificate, []byte, error) {
	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, subjectKey, issuerKey)
	if err != nil {
		return nil, nil, err
	}
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, nil, err
	}

	b := pem.Block{Type: "CERTIFICATE", Bytes: certBytes}
	certPEM := pem.EncodeToMemory(&b)

	return cert, certPEM, nil
}

func GetRandomSerial() *big.Int {
	z := new(big.Int)
	b, err := genRandomBytes(4)
	if err != nil {
		log.Fatalf("Failed to generate random serial, returned error: %s\n", err)
	}
	z.SetBytes(b)
	return z
}

func (cert *Certificate) loadCertAndKey(filename string) error {
	privateKey, c, cChain, err := ReadPFX(filename)
	if err != nil {
		return err
	}

	cert.Cert = c
	cert.PrivateKey = privateKey
	cert.CertChain = cChain
	cert.Type = GetType(cert.Cert)

	cert.CertPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: c.Raw,
	})
	cert.PrivateKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	return nil
}

func genRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, errors.New("failed generating random bytes, length less than 0")
	}
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
