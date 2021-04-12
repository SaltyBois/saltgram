package handlers

import (
	"net/http"
	"saltgram/data"

	"github.com/gorilla/mux"
)

func (e *Emails) Activate(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	token := vars["token"]

	err := data.ActivateEmail(token)
	if err != nil {
		e.l.Printf("[ERROR] activating email: %v", err)
		http.Error(w, "Failed activating email: " + err.Error(), http.StatusInternalServerError)
		return
	}
	data.ToJSON("Email activated!", w)
}

func (e *Emails) GetAll(w http.ResponseWriter, r *http.Request) {
	err := data.ToJSON(data.GetAllActivations(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (u *Users) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id, err := getUserID(r)
	if err != nil {
		http.Error(w, "Invalid User ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := data.GetUserByID(id)

	switch err {
	case nil:
	case data.ErrUserNotFound:
		u.l.Println("[ERROR] fetching user", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	default:
		u.l.Println("[ERROR] fetching user", err)
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	err = data.ToJSON(user, w)
	if err != nil {
		u.l.Println("[ERROR] serializing user ", err)
		http.Error(w, "Error serializing user", http.StatusInternalServerError)
		return
	}

}

func (u *Users) GetAll(w http.ResponseWriter, r *http.Request) {
	u.l.Println("Handling GET Users")

	ul := data.GetUsers()
	err := data.ToJSON(ul, w)
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}
