package internal

import (
	"crypto/tls"
	"crypto/x509"
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

func (s *Service) Init(serverName string, certPEMBlock, keyPEMBlock, rootPEMBlock []byte) error {
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return err
	}
	conf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   serverName,
		MinVersion:   tls.VersionTLS13,
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(rootPEMBlock)
	conf.RootCAs = caPool
	s.TLS.TC = conf
	s.TLS.C = credentials.NewTLS(s.TLS.TC)
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
