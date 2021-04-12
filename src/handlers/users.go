package handlers

import (
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

func getUserJWS(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	// NOTE(Jovan): Trimming first 7 characters from "Bearer <jws>"
	return authHeader[7:]

	// vars := mux.Vars(r)
	// return vars["jws"]
}
