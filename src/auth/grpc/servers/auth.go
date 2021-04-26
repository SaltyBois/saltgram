package servers

import (
	"context"
	"fmt"
	"log"
	"os"
	"saltgram/auth/data"
	saltdata "saltgram/data"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/users/prusers"
	"strconv"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
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

func (a *Auth) Refresh(ctx context.Context, r *prauth.RefreshRequest) (*prauth.RefreshResponse, error) {

	rTokenString := r.Refresh
	rToken, err := jwt.ParseWithClaims(
		rTokenString,
		&saltdata.RefreshClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("REF_SECRET_KEY")), nil
		},
	)

	if err != nil {
		a.l.Printf("[ERROR] parsing refresh claims: %v", err)
		return &prauth.RefreshResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	claims, ok := rToken.Claims.(*saltdata.RefreshClaims)

	if !ok {
		a.l.Println("[ERROR] unable to parse claims")
		return &prauth.RefreshResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	refreshToken, err := data.GetRefreshToken(a.db, claims.Username)

	if err != nil {
		a.l.Println("[ERROR] can't find refresh token")
		return &prauth.RefreshResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	rt := data.Refresh{
		Username: claims.Username,
		Token:    refreshToken,
	}
	if err := rt.Verify(a.db); err != nil {
		a.l.Println("[ERROR] refresh token no longer valid")
		return &prauth.RefreshResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	jws := r.OldJWS

	// NOTE(Jovan): Not validating 'cause it is invalid
	jwtOld, _ := jwt.ParseWithClaims(
		jws,
		&saltdata.AccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		},
	)

	jwsClaims, ok := jwtOld.Claims.(*saltdata.AccessClaims)

	if !ok {
		a.l.Println("[ERROR] unable to parse claims")
		return &prauth.RefreshResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	// TODO(Jovan): Pull out into const
	timeout, _ := strconv.Atoi(os.Getenv("TOKEN_TIMEOUT_MINUTES"))
	jwsClaims.StandardClaims.ExpiresAt = time.Now().UTC().Add(time.Minute * time.Duration(timeout)).Unix()
	jwtNew := jwt.NewWithClaims(jwt.SigningMethodHS256, jwsClaims)

	jwsNew, err := jwtNew.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		a.l.Printf("[ERROR] failed signing JWT: %v", err)
		return &prauth.RefreshResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prauth.RefreshResponse{NewJWS: jwsNew}, nil
}

func (a *Auth) GetJWT(ctx context.Context, r *prauth.JWTRequest) (*prauth.JWTResponse, error) {

	// NOTE(Jovan): HS256 is considered safe enough
	claims := saltdata.AccessClaims{
		Username: r.Username,
		Password: r.Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Second * 5).Unix(),
			Issuer:    "SaltGram",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jws, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		a.l.Printf("[ERROR] failed signing JWT: %v", err)
		return &prauth.JWTResponse{}, status.Error(codes.Internal, "Internal server error")
	}

	refreshToken, err := data.GetRefreshToken(a.db, r.Username)
	if err != nil {
		a.l.Printf("[ERROR] failed getting refresh token: %v", err)
		return &prauth.JWTResponse{}, status.Error(codes.Internal, "Internal server error")
	}

	return &prauth.JWTResponse{
		Jws:     jws,
		Refresh: refreshToken,
	}, nil
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

	pres, err := a.usersClient.CheckPassword(context.Background(), &prusers.CheckPasswordRequest{Username: r.Username, Password: r.Password})
	if err != nil {
		a.l.Println("[ERROR] bad request")
		return &prauth.LoginResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prauth.LoginResponse{Username: r.Username, Password: pres.HashedPass}, nil
}
