package internal

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Service struct {
	L   *log.Logger
	S   *mux.Router
	TLS *TLS
}

type TLS struct {
	TC *tls.Config
	C  credentials.TransportCredentials
}

func GetEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func GetEnvOrDefaultInt(key string, fallback int) int {
	value, err := strconv.ParseInt(os.Getenv(key), 10, 32)
	if err != nil {
		return fallback
	}
	return int(value)
}

func NewService(l *log.Logger) *Service {
	return &Service{L: l, S: mux.NewRouter(), TLS: &TLS{}}
}

func (t *TLS) Init(localCrtPath, localKeyPath, rootPEMPath string) error {

	cert, err := loadTLSCert(localCrtPath, localKeyPath)
	if err != nil {
		return err
	}

	conf := &tls.Config{
		Certificates: []tls.Certificate{*cert},
		ServerName:   "localhost",
		MinVersion:   tls.VersionTLS13,
	}
	if len(rootPEMPath) > 0 {
		caCert, err := ioutil.ReadFile(rootPEMPath)
		if err != nil {
			return err
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCert)
		conf.RootCAs = caPool
	}
	t.TC = conf
	t.C = credentials.NewTLS(t.TC)
	return nil
}

func (s *Service) NewServer() *grpc.Server {
	return grpc.NewServer(grpc.Creds(s.TLS.C))
}

func (s *Service) GetConnection(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(s.TLS.C))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func loadTLSCert(crtPath, keyPath string) (*tls.Certificate, error) {
	crt, err := ioutil.ReadFile(crtPath)
	if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}
