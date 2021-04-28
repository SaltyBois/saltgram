package handlers

import (
	"context"
	"net/http"
	"saltgram/auth/data"
	saltdata "saltgram/data"
)

// func (l Login) MiddlewareValidateToken(e *casbin.Enforcer) func(next http.Handler) http.Handler {
func (l Login) MiddlewareValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login := saltdata.Login{}

		err := saltdata.FromJSON(&login, r.Body)
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

func (a *Auth) RefreshMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := data.Refresh{}

		// NOTE(Jovan): Deserialize JSON object
		err := saltdata.FromJSON(&token, r.Body)
		if err != nil {
			a.l.Println("[ERROR] deserializing refresh token: ", err.Error())
			http.Error(w, "Error reading user", http.StatusBadRequest)
			return
		}

		// NOTE(Jovan): Validate token
		err = token.Validate()
		if err != nil {
			a.l.Println("[ERROR] validating token: ", err.Error())
			http.Error(w, "Error validating token: "+err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyRefreshToken{}, token)
		requestCopy := r.WithContext(ctx)
		next.ServeHTTP(w, requestCopy)
	})
}

// func (a *Auth) AuthMiddleware(e *casbin.Enforcer) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 	}
// }
