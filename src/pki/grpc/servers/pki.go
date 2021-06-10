package servers

import (
	"context"
	"crypto/x509/pkix"
	"log"
	"saltgram/pki/data"
	"saltgram/protos/pki/prpki"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PKI struct {
	prpki.UnimplementedPKIServer
	l  *log.Logger
	db *data.DBConn
}

func NewPKI(l *log.Logger, db *data.DBConn) *PKI {
	return &PKI{l: l, db: db}
}

func (p *PKI) RegisterService(ctx context.Context, r *prpki.PKIRegisterRequest) (*prpki.PKIRegisterResponse, error) {
	_, err := data.RegisterService(p.db, pkix.Name{
		CommonName: r.CommonName,
		Organization: r.Organization,
		Country: r.Country,
		Province: r.Province,
		Locality: r.Locality,
		StreetAddress: r.StreetAddress,
		PostalCode: r.PostalCode,
	})
	if err != nil {
		p.l.Printf("[ERROR] registering service: %v\n", err)
		return &prpki.PKIRegisterResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}
	return &prpki.PKIRegisterResponse{}, nil
}