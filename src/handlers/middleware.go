package handlers

import (
	"context"
	"net/http"
	"saltgram/data"
)

func (re Login) MiddlewareValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(Jovan): Move
		w.Header().Add("Strict-Transport-Security", "max-age=86400")
		reCaptcha := data.Login{}

		err := data.FromJSON(&reCaptcha, r.Body)
		if err != nil {
			re.l.Println("[ERROR] deserializing reCaptcha token: ", err.Error())
			http.Error(w, "Error getting reCaptcha token", http.StatusBadRequest)
			return
		}

		err = reCaptcha.Validate()
		if err != nil {
			re.l.Println("[ERROR] validating reCaptcha token: ", err.Error())
			http.Error(w, "Error validating reCaptcha token", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyLogin{}, reCaptcha)
		requestCopy := r.WithContext(ctx)

		next.ServeHTTP(w, requestCopy)
	})
}

func (u Users) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(Jovan): Move
		w.Header().Add("Strict-Transport-Security", "max-age=86400")
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

		// NOTE(Jovan): If JSON object is valid, put the unmarshalled
		// struct onto request
		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		requestCopy := r.WithContext(ctx)

		next.ServeHTTP(w, requestCopy)
	})
}
