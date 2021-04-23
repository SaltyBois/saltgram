package handlers

import (
	"context"
	"net/http"
	"saltgram/auth/data"
)

func (a *Auth) RefreshMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := data.Refresh{}

		// NOTE(Jovan): Deserialize JSON object
		err := data.FromJSON(&token, r.Body)
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