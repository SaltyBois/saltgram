package servers

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/users/prusers"
	"saltgram/users/data"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Users struct {
	prusers.UnimplementedUsersServer
	l  *log.Logger
	db *data.DBConn
	ac prauth.AuthClient
}

func NewUsers(l *log.Logger, db *data.DBConn, ac prauth.AuthClient) *Users {
	return &Users{
		l:  l,
		db: db,
		ac: ac,
	}
}

func (u *Users) Register(ctx context.Context, r *prusers.RegisterRequest) (*prusers.RegisterResponse, error) {
	u.l.Println("Handling POST Users")

	refreshClaims := data.RefreshClaims{
		Username: r.Username,
		StandardClaims: jwt.StandardClaims{
			// TODO(Jovan): Make programmatic?
			ExpiresAt: time.Now().UTC().AddDate(0, 6, 0).Unix(),
			Issuer:    "SaltGram",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	jws, err := token.SignedString([]byte(os.Getenv("REF_SECRET_KEY")))

	if err != nil {
		u.l.Println("[ERROR] signing refresh token")
		return &prusers.RegisterResponse{}, status.Error(codes.Internal, "Internal server error")
	}

	user := data.User{
		Email:          r.Email,
		Username:       r.Username,
		FullName:       r.FullName,
		HashedPassword: r.Password,
		Role:           "user", // TODO(Jovan): For now
	}
	err = u.db.AddUser(&user)
	if err != nil {
		u.l.Printf("[ERROR] adding user: %v\n", err)
		return &prusers.RegisterResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	// NOTE(Jovan): Saving refresh token
	// data.AddRefreshToken(user.Username, jws)

	// resp, err := http.Post("https://localhost:8082/refresh", "application/json", bytes.NewBuffer(jsonData))
	_, err = u.ac.AddRefresh(context.Background(), &prauth.AddRefreshRequest{
		Username: r.Username,
		Token:    jws,
	})

	if err != nil {
		u.l.Printf("[ERROR] adding refresh token: %v\n", err)
		return &prusers.RegisterResponse{}, status.Error(codes.Internal, "Internal server error")
	}

	go func() {
		resp, err := http.Post("https://localhost:8084/activate", "text/html", bytes.NewBuffer([]byte(user.Email)))
		if err != nil {
			u.l.Printf("[ERROR] sending change password request: %v\n", err)
		}

		if resp.StatusCode != http.StatusOK {
			u.l.Println("[ERROR] sending change password request")
		}
	}()

	return &prusers.RegisterResponse{}, nil
}

func (u *Users) CheckPassword(ctx context.Context, r *prusers.CheckPasswordRequest) (*prusers.CheckPasswordResponse, error) {
	hashedPass, err := data.IsPasswordValid(u.db, r.Username, r.Password)
	if err != nil {
		u.l.Printf("[ERROR] checking password: %v\n", err)
		return &prusers.CheckPasswordResponse{}, err
	}
	return &prusers.CheckPasswordResponse{HashedPass: hashedPass}, nil
}

func (u *Users) CheckEmail(ctx context.Context, r *prusers.CheckEmailRequest) (*prusers.CheckEmailResponse, error) {
	return &prusers.CheckEmailResponse{
		Verified: data.IsEmailVerified(u.db, r.Username),
	}, nil
}
