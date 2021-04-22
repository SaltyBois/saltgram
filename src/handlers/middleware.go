package handlers

import (
	"context"
	"net/http"
	"saltgram/data"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func MultipleMiddleware(h http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}

	wrapped := h
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}
	return wrapped
}

// func (l Login) MiddlewareValidateToken(e *casbin.Enforcer) func(next http.Handler) http.Handler {
func (l Login) MiddlewareValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login := data.Login{}

		err := data.FromJSON(&login, r.Body)
		if err != nil {
			l.l.Println("[ERROR] deserializing reCaptcha token: ", err.Error())
			http.Error(w, "Error getting reCaptcha token", http.StatusBadRequest)
			return
		}

		err = login.Validate()
		if err != nil {
			l.l.Println("[ERROR] validating reCaptcha token: ", err.Error())
			http.Error(w, "Error validating reCaptcha token", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyLogin{}, login)
		requestCopy := r.WithContext(ctx)

		next.ServeHTTP(w, requestCopy)
	})
}

