package handlers

import (
	"net/http"

	"github.com/casbin/casbin/v2"
)

func AuthMiddleware(e *casbin.Enforcer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}