package handlers

import (
	"log"
	"net/http"
	"saltgram/data"
)

type Users struct {
	l *log.Logger
}

func NewUsers(l *log.Logger) *Users {
	return &Users{l}
}

func (u *Users) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// NOTE(Jovan): Handle GET method
	if r.Method == http.MethodGet {
		u.getUsers(w, r)
		return
	}

	// NOTE(Jovan): If no method is satisfied, catchall and return err
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (u *Users) getUsers(w http.ResponseWriter, r *http.Request) {
	u.l.Println("Handling GET users")

	ul := data.GetUsers()

	err := ul.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}
