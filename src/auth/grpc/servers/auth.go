package servers

import (
	"context"
	"fmt"
	"log"
	"saltgram/auth/data"
	saltdata "saltgram/data"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/users/prusers"

	"github.com/casbin/casbin/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth struct {
	prauth.UnimplementedAuthServer
	l           *log.Logger
	e           *casbin.Enforcer
	db          *data.DBConn
	usersClient prusers.UsersClient
}

func NewAuth(l *log.Logger, e *casbin.Enforcer, db *data.DBConn, usersClient prusers.UsersClient) *Auth {
	return &Auth{
		l:           l,
		e:           e,
		db:          db,
		usersClient: usersClient,
	}
}

func (a *Auth) AddRefresh(ctx context.Context, r *prauth.AddRefreshRequest) (*prauth.AddRefreshResponse, error) {
	token := data.Refresh{
		Username: r.Username,
		Token:    r.Token,
	}
	err := data.AddRefreshToken(a.db, token.Username, token.Token)
	if err != nil {
		a.l.Printf("[ERROR] adding refresh token: %v\n", err)
		return &prauth.AddRefreshResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prauth.AddRefreshResponse{}, nil
}

var ErrorBadRequest = fmt.Errorf("bad request")

func (a *Auth) Login(ctx context.Context, r *prauth.LoginRequest) (*prauth.LoginResponse, error) {
	a.l.Print("Doing stuff")
	res, err := a.usersClient.CheckEmail(context.Background(), &prusers.CheckEmailRequest{Username: r.Username})
	if err != nil {
		a.l.Printf("[ERROR] checking email: %v\n", err)
		return &prauth.LoginResponse{}, err
	}

	recaptcha := saltdata.ReCaptcha{
		Action: r.ReCaptcha.Action,
		Token:  r.ReCaptcha.Token,
	}
	score, err := recaptcha.Verify()
	if err != nil {
		a.l.Println("[ERROR] verifying reCaptcha")
		return &prauth.LoginResponse{}, err
	}

	if score < 0.5 {
		return &prauth.LoginResponse{}, status.Error(codes.PermissionDenied, "Low ReCaptcha score")
	}

	if !res.Verified {
		a.l.Println("[ERROR] bad request")
		return &prauth.LoginResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	pres, err := a.usersClient.CheckPassword(context.Background(), &prusers.CheckPasswordRequest{Password: r.Password})
	if err != nil {
		a.l.Println("[ERROR] bad request")
		return &prauth.LoginResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prauth.LoginResponse{Username: r.Username, Password: pres.HashedPass}, nil
}
