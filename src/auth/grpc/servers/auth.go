package servers

import (
	"context"
	"fmt"
	"os"
	"saltgram/auth/data"
	saltdata "saltgram/data"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/users/prusers"
	"strconv"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth struct {
	prauth.UnimplementedAuthServer
	l  *logrus.Logger
	e  *casbin.Enforcer
	db *data.DBConn
	uc prusers.UsersClient
}

func NewAuth(l *logrus.Logger, e *casbin.Enforcer, db *data.DBConn, uc prusers.UsersClient) *Auth {
	return &Auth{
		l:  l,
		e:  e,
		db: db,
		uc: uc,
	}
}

func (a *Auth) Authenticate2FA(ctx context.Context, r *prauth.Auth2FARequest) (*prauth.Auth2FAResponse, error) {
	a.l.Infof("authenticating 2FA with token: %v\n", r.Token)
	ok, err := data.Authenticate2FA(r.Token)
	if err != nil {
		a.l.Errorf("failure authenticating 2FA token: %v, error: %v\n", r.Token, err)
		return &prauth.Auth2FAResponse{}, status.Error(codes.Internal, "Internal server error")
	}

	if !ok {
		a.l.Warnf("failure authenticating 2FA token: %v\n", r.Token)
		return &prauth.Auth2FAResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}
	return &prauth.Auth2FAResponse{}, nil
}

func (a *Auth) Get2FAQR(ctx context.Context, r *prauth.TwoFARequest) (*prauth.TwoFAResponse, error) {
	a.l.Infof("doing 2FA for: %v\n", r.Username)
	res, err := a.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: r.Username})
	if err != nil {
		a.l.Errorf("failed to get user by username: %v, error %v\n", r.Username, err)
		return &prauth.TwoFAResponse{Png: nil}, status.Error(codes.InvalidArgument, "Bad request")
	}
	png, err := data.Get2FAQR(res.Email)
	if err != nil {
		a.l.Errorf("failed to create 2FA QR: %v\n", err)
		return &prauth.TwoFAResponse{Png: nil}, status.Error(codes.Internal, "Internal server error")
	}
	return &prauth.TwoFAResponse{Png: png}, nil
}

func (a *Auth) CheckPermissions(ctx context.Context, r *prauth.PermissionRequest) (*prauth.PermissionResponse, error) {

	subject := "anon"
	if len(r.Jws) > 0 {
		jwt, err := jwt.ParseWithClaims(
			r.Jws,
			&saltdata.AccessClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			},
		)
		if err != nil {
			a.l.Errorf("failure parsing token: %v\n", err)
			return &prauth.PermissionResponse{}, status.Error(codes.InvalidArgument, "Bad request")
		}
		claims, ok := jwt.Claims.(*saltdata.AccessClaims)
		if !ok {
			a.l.Errorf("failure parsing claims: %v\n", err)
			return &prauth.PermissionResponse{}, status.Error(codes.InvalidArgument, "Bad request")
		}
		res, err := a.uc.GetRole(context.Background(), &prusers.RoleRequest{Username: claims.Username})
		if err != nil {
			a.l.Errorf("failed to get user role: %v\n", err)
			return &prauth.PermissionResponse{}, status.Error(codes.InvalidArgument, "Bad request")
		}
		subject = res.Role
	}
	ok, err := a.e.Enforce(subject, r.Path, r.Method)
	if err != nil {
		a.l.Errorf("failure enforcing: %v\n", err)
		return &prauth.PermissionResponse{}, status.Error(codes.Internal, "Internal server error")
	}
	if !ok {
		a.l.Infof("[DENIED] Subject: %v, Object: %v, Method: %v\n", subject, r.Path, r.Method)
		return &prauth.PermissionResponse{}, status.Error(codes.PermissionDenied, "Permission denied")
	}
	a.l.Infof("[GRANTED] Subject: %v, Object: %v, Method: %v\n", subject, r.Path, r.Method)
	return &prauth.PermissionResponse{}, nil
}

func (a *Auth) Refresh(ctx context.Context, r *prauth.RefreshRequest) (*prauth.RefreshResponse, error) {
	rToken, err := jwt.ParseWithClaims(
		r.Refresh,
		&saltdata.RefreshClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("REF_SECRET_KEY")), nil
		},
	)

	if err != nil {
		a.l.Errorf("failure parsing refresh claims: %v\n", err)
		return &prauth.RefreshResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	claims, ok := rToken.Claims.(*saltdata.RefreshClaims)

	if !ok {
		a.l.Error("unable to parse claims")
		return &prauth.RefreshResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	refreshToken, err := data.GetRefreshToken(a.db, claims.Username)

	if err != nil {
		a.l.Errorf("can't find refresh token: %v\n")
		return &prauth.RefreshResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	rt := data.Refresh{
		Username: claims.Username,
		Token:    refreshToken,
	}
	if err := rt.Verify(a.db); err != nil {
		a.l.Errorf("refresh token: %v no longer valid: %v\n", rt.Token, err)
		return &prauth.RefreshResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	// NOTE(Jovan): Not validating 'cause it is invalid
	jwtOld, _ := jwt.ParseWithClaims(
		r.OldJWS,
		&saltdata.AccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("REF_SECRET_KEY")), nil
		},
	)

	jwsClaims, ok := jwtOld.Claims.(*saltdata.AccessClaims)

	if !ok {
		a.l.Error("unable to parse claims")
		return &prauth.RefreshResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	// TODO(Jovan): Pull out into const
	timeout, _ := strconv.Atoi(os.Getenv("TOKEN_TIMEOUT_MINUTES"))
	jwsClaims.StandardClaims.ExpiresAt = time.Now().UTC().Add(time.Minute * time.Duration(timeout)).Unix()
	jwtNew := jwt.NewWithClaims(jwt.SigningMethodHS256, jwsClaims)

	jwsNew, err := jwtNew.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		a.l.Errorf("failed signing JWT: %v\n", err)
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
		a.l.Errorf("failed signing JWT: %v\n", err)
		return &prauth.JWTResponse{}, status.Error(codes.Internal, "Internal server error")
	}

	refreshToken, err := data.GetRefreshToken(a.db, r.Username)
	if err != nil {
		a.l.Errorf("failed getting refresh token: %v\n", err)
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
		a.l.Errorf("adding refresh token: %v\n", err)
		return &prauth.AddRefreshResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prauth.AddRefreshResponse{}, nil
}

var ErrorBadRequest = fmt.Errorf("bad request")

func (a *Auth) Login(ctx context.Context, r *prauth.LoginRequest) (*prauth.LoginResponse, error) {
	res, err := a.uc.CheckEmail(context.Background(), &prusers.CheckEmailRequest{Username: r.Username})
	if err != nil {
		a.l.Errorf("checking email: %v\n", err)
		return &prauth.LoginResponse{}, err
	}

	recaptcha := saltdata.ReCaptcha{
		Action: r.ReCaptcha.Action,
		Token:  r.ReCaptcha.Token,
	}
	score, err := recaptcha.Verify()
	if err != nil {
		a.l.Errorf("failed to verify reCaptcha: %v\n", err)
		return &prauth.LoginResponse{}, err
	}

	if score < 0.5 {
		a.l.Infof("[DENIED] due to low recaptcha score: %v\n", score)
		return &prauth.LoginResponse{}, status.Error(codes.PermissionDenied, "Low ReCaptcha score")
	}

	if !res.Verified {
		a.l.Error("mail not verified")
		return &prauth.LoginResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	pres, err := a.uc.CheckPassword(context.Background(), &prusers.CheckPasswordRequest{Username: r.Username, Password: r.Password})
	if err != nil {
        a.l.Errorf("failure validating password: %v\n", err)
		return &prauth.LoginResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prauth.LoginResponse{Username: r.Username, Password: pres.HashedPass}, nil
}
