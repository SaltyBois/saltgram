package data

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"saltgram/internal"

	"software.sslmate.com/src/go-pkcs12"
)

var (
	ROOT_DIR = filepath.FromSlash(internal.GetEnvOrDefault("KEYSTORE_ROOT_DIR", "../keystore/"))
	FILE_EXT = ".pfx"
)

func InitKeystore() {
	if _, err := os.Stat(ROOT_DIR); os.IsNotExist(err) {
		os.Mkdir(ROOT_DIR, 0755)
	}
}

func ReadPFX(filename string) (*rsa.PrivateKey, *x509.Certificate, []*x509.Certificate, error) {
	pfxData, err := ioutil.ReadFile(filepath.FromSlash(ROOT_DIR + filename + FILE_EXT))
	if err != nil {
		return nil, nil, nil, err
	}

	privateKey, cert, caCerts, err := pkcs12.DecodeChain(pfxData, pkcs12.DefaultPassword)
	if err != nil {
		return nil, nil, nil, err
	}

	PKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, nil, errors.New("could not convert to rsa.PrivateKey")
	}
	return PKey, cert, caCerts, nil
}

func WritePFX(cert *x509.Certificate, certChain []*x509.Certificate, PrivateKey *rsa.PrivateKey, filename string) error {
	pfxBytes, err := pkcs12.Encode(rand.Reader, PrivateKey, cert, certChain, pkcs12.DefaultPassword)
	if err != nil {
		return err
	}
	if err := os.WriteFile(
		ROOT_DIR+filename+FILE_EXT,
		pfxBytes,
		os.ModePerm,
	); err != nil {
		return err
	}
	return nil
}
