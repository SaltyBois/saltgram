package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"saltgram/pki"
	"time"

	"saltgram/log"

	spa "github.com/roberthodgen/spa-server"
)

func getTLSConfig(certPEM, keyPEM []byte, pki *pki.PKI) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	rootPool := x509.NewCertPool()
	rootPool.AddCert(pki.RootCA.Cert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
		MinVersion:   tls.VersionTLS13,
		// RootCAs: rootPool,
	}, nil
}

func hstsMiddleware(h http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=86400")
		h.ServeHTTP(w, r)
	}
}

func main() {
	l := log.NewLogger("saltgram-webserver")
	l.L.Infof("Starting webserver on port: %s", os.Getenv("SALT_WEB_PORT"))
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", hstsMiddleware(spa.SpaHandler("../frontend/dist", "index.html")))
	pkiHandler := pki.Init()
	cert, err := pkiHandler.RegisterSaltgramService("saltgram-web-server")
	if err != nil {
		l.L.Fatalf("[ERROR] registering to PKI: %v\n", err)
	}
	tlsConfig, err := getTLSConfig(cert.CertPEM, cert.PrivateKeyPEM, pkiHandler)
	if err != nil {
		l.L.Fatalf("failed to get TLS config: %v\n", err)
	}

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("SALT_WEB_PORT")),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      serverMux,
	}

	go func() {
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			l.L.Fatalf("[ERROR] while serving: %v\n", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	sig := <-signalChan
	l.L.Infof("Recieved terminate, graceful shutdown with sigtype: %v\n", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(tc)
}
