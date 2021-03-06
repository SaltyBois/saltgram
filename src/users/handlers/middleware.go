package handlers

import (
	"context"
	"net/http"
	"saltgram/users/data"
)

func (u Users) MiddlewareValidateChangeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			next.ServeHTTP(w, r)
			return
		}
		changeRequest := ChangeRequest{}

		err := data.FromJSON(&changeRequest, r.Body)
		if err != nil {
			u.l.Printf("[ERROR] deserializing password change request: %v\n", err)
			http.Error(w, "Bad password change request", http.StatusBadRequest)
			return
		}

		err = changeRequest.Validate()
		if err != nil {
			u.l.Printf("[ERROR] validating password change request: %v\n", err)
			http.Error(w, "Invalid password change request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUser{}, changeRequest)
		requestCopy := r.WithContext(ctx)

		next.ServeHTTP(w, requestCopy)
	})
}

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

		// NOTE(Jovan): Validate user
		err = user.Validate()
		if err != nil {
			u.l.Println("[ERROR] validating user: ", err.Error())
			http.Error(w, "Error validating user: "+err.Error(), http.StatusBadRequest)
			return
		}

		// NOTE(Jovan): Role -> user when registering
		// TODO(Jovan): Un-hardcode, tidy it up
		if r.Method == http.MethodPost {
			user.Role = "user"
		}
		// else {
		// 	role, err := data.GetRole(user.Username)
		// 	if err != nil {
		// 		u.l.Printf("[ERROR] getting user role: %v\n", err)
		// 		http.Error(w, "Invalid user role", http.StatusBadRequest)
		// 		return
		// 	}
		// 	u.l.Printf("Requested path: %v for method: %v as role: %v\n", r.URL.Path, r.Method, role)
		// 	res, err := e.Enforce(role, r.URL.Path, r.Method)
		// 	if err != nil {
		// 		u.l.Printf("[ERROR] while enforcing: %v\n", err)
		// 		http.Error(w, "Error while enforcing", http.StatusInternalServerError)
		// 		return
		// 	}

		// 	if !res {
		// 		u.l.Printf("Forbidden access! Subject: %v, Object: %v, Act: %v\n", role, r.URL.Path, r.Method)
		// 		http.Error(w, "Forbidden access", http.StatusForbidden)
		// 		return
		// 	}
		// }

		// NOTE(Jovan): If JSON object is valid, put the unmarshalled
		// struct onto request
		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		requestCopy := r.WithContext(ctx)

		next.ServeHTTP(w, requestCopy)
	})
}
