package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Users struct {
	l *log.Logger
}

// NOTE(Jovan): Key used for contexts
type KeyUser struct{}

func NewUsers(l *log.Logger) *Users {
	return &Users{l}
}

func getUserID(r *http.Request) (uint64, error) {
	vars := mux.Vars(r)
	idstring := vars["id"]

	id, err := strconv.ParseUint(idstring, 10, 64)
	return id, err
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
