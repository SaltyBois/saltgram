package handlers

import (
	"net/http"
	"saltgram/users/data"
)

func (u *Users) Update(db *data.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getUserID(r)

		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		u.l.Println("Handling PUT Users", id)
		// NOTE(Jovan): Safe to cast because middleware makes sure nothing's wrong
		user := r.Context().Value(KeyUser{}).(data.User)

		err = db.UpdateUser(&user)
		if err != nil {
			u.l.Printf("[ERROR] updating user: %v\n", err)
			http.Error(w, "Failed to update user", http.StatusNotFound)
			return
		}
	}
}
