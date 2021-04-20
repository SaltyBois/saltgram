package handlers

import (
	"context"
	"net/http"
	"saltgram/data"

	"github.com/casbin/casbin/v2"
)

func (l Login) MiddlewareValidateToken(e *casbin.Enforcer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			role, err := data.GetRole(login.Username)
			if err != nil {
				l.l.Printf("[ERROR] getting user role: %v\n", err)
				http.Error(w, "Invalid user role", http.StatusBadRequest)
				return
			}
			l.l.Printf("Requested path: %v for method: %v as role: %v\n", r.URL.Path, r.Method, role)

			ctx := context.WithValue(r.Context(), KeyLogin{}, login)
			requestCopy := r.WithContext(ctx)

			next.ServeHTTP(w, requestCopy)
		})
	}
}

func (u Users) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := data.User{}

		// NOTE(Jovan): Validate JSON object
		err := data.FromJSON(&user, r.Body)
		if err != nil {
			u.l.Println("[ERROR] deserializing user aaa: ", err.Error())
			http.Error(w, "Error reading user", http.StatusBadRequest)
			return
		}

		// NOTE(Jovan): Validate product
		err = user.Validate()
		if err != nil {
			u.l.Println("[ERROR] validating user: ", err.Error())
			http.Error(w, "Error validating product: "+err.Error(), http.StatusBadRequest)
			return
		}

		// NOTE(Jovan): Role -> user
		// TODO(Jovan): Un-hardcode, tidy it up
		user.Role = "user"

		// NOTE(Jovan): If JSON object is valid, put the unmarshalled
		// struct onto request
		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		requestCopy := r.WithContext(ctx)

		next.ServeHTTP(w, requestCopy)
	})
}
