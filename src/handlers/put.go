package handlers

import (
	"net/http"
	"saltgram/data"
)


func (u *Users) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getUserID(r)

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
