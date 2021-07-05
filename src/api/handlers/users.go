package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	saltdata "saltgram/data"
	"saltgram/protos/users/prusers"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type Users struct {
	l  *logrus.Logger
	uc prusers.UsersClient
}

func NewUsers(l *logrus.Logger, uc prusers.UsersClient) *Users {
	return &Users{l: l, uc: uc}
}

var ErrorJWSNotFound = fmt.Errorf("jws not found")

func getUserJWS(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) <= 7 {
		return "", ErrorJWSNotFound
	}
	// NOTE(Jovan): Trimming first 7 characters from "Bearer <jws>"
	return authHeader[7:], nil
}

func getUserByJWS(r *http.Request, uc prusers.UsersClient) (*prusers.GetByUsernameResponse, error) {
	jws, err := getUserJWS(r)
	if err != nil {
		return nil, fmt.Errorf("JWS not found: %v", err)
	}

	token, err := jwt.ParseWithClaims(
		jws,
		&saltdata.AccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failure parsing claims: %v", err)
	}
	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}

	return uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: claims.Username})
}

func getProfileByJWS(r *http.Request, uc prusers.UsersClient) (*prusers.ProfileResponse, error) {
	jws, err := getUserJWS(r)
	if err != nil {
		return nil, fmt.Errorf("JWS not found: %v", err)
	}

	token, err := jwt.ParseWithClaims(
		jws,
		&saltdata.AccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed parsing claims: %v", err)
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		return nil, fmt.Errorf("unable to parse claims")
	}

	fmt.Println("claims: ", claims.Username)

	return uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{Username: claims.Username, User: claims.Username})
}

func getUsernameByJWS(r *http.Request) (string, error) {
	jws, err := getUserJWS(r)
	if err != nil {
		return "", fmt.Errorf("JWS not found: %v", err)
	}

	token, err := jwt.ParseWithClaims(
		jws,
		&saltdata.AccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed parsing claims: %v", err)
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		return "", fmt.Errorf("unable to parse claims")
	}

	return claims.Username, nil

}
