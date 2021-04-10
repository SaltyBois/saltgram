package handlers

import (
	"net/http"
	"saltgram/data"
)

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	u.l.Println("Handling POST Users")

	user := r.Context().Value(KeyUser{}).(data.User)
	data.AddUser(&user)
}