package handlers

import (
	"log"
	"net/http"
	"saltgram/data"
	"strconv"

	"github.com/gorilla/mux"
)

type Users struct {
	l *log.Logger
}

type KeyUser struct {}

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

	user := &data.User{}

	err = user.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	err = data.UpdateUser(id, user)
	if err == data.ErrUserNotFound {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Product not updated: " + err.Error(), http.StatusInternalServerError)
		return
	}
}
