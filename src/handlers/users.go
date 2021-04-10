package handlers

import (
	"context"
	"log"
	"net/http"
	"saltgram/data"
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

func (u *Users) GetUsers(w http.ResponseWriter, r *http.Request) {
	u.l.Println("Handling GET Users")

	ul := data.GetUsers()

	err := ul.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}

func (u *Users) AddUser(w http.ResponseWriter, r *http.Request) {
	u.l.Println("Handling POST Users")

	user := r.Context().Value(KeyUser{}).(data.User)
	data.AddUser(&user)
}

func (u *Users) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstring := vars["id"]

	id, err := strconv.ParseUint(idstring, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	u.l.Println("Handling PUT Users", id)
	// NOTE(Jovan): Safe to cast because middleware makes sure nothing's wrong
	user := r.Context().Value(KeyUser{}).(data.User)

	err = data.UpdateUser(id, &user)
	if err == data.ErrUserNotFound {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Product not updated: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u Users) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := data.User{}

		// NOTE(Jovan): Validate JSON object
		err := user.FromJSON(r.Body)
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
