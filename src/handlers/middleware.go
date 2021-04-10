package handlers

import (
	"context"
	"net/http"
	"saltgram/data"
)


func (u Users) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := data.User{}

		// NOTE(Jovan): Validate JSON object
		err := data.FromJSON(&user, r.Body)
		if err != nil {
			u.l.Println("[ERROR] deserializing user: ", err.Error())
			http.Error(w, "Error reading user", http.StatusBadRequest)
			return
		}

		// NOTE(Jovan): Validate product
		err = user.Validate()
		if err != nil {
			u.l.Println("[ERROR] validating user: ", err.Error())
			http.Error(w, "Error validating product: " + err.Error(), http.StatusBadRequest)
			return
		}

		// NOTE(Jovan): If JSON object is valid, put the unmarshalled
		// struct onto request
		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		requestCopy := r.WithContext(ctx)

		next.ServeHTTP(w, requestCopy)
	})
}