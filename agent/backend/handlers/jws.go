package handlers

import (
	"fmt"
	"net/http"
)

var ErrorJWSNotFound = fmt.Errorf("jws not found")

func getUserJWS(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) <= 7 {
		return "", ErrorJWSNotFound
	}
	// NOTE(Jovan): Trimming first 7 characters from "Bearer <jws>"
	return authHeader[7:], nil
}
