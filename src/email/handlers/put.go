package handlers

import (
	"bytes"
	"net/http"
	"saltgram/email/data"
	"time"

	"github.com/gorilla/mux"
)


func (e *Emails) ConfirmReset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	email, err := data.ConfirmPasswordReset(token)
	if err != nil {
		e.l.Printf("[ERROR] confirming password reset: %v\n", err)
		http.Error(w, "Failed to confirm password reset", http.StatusBadRequest)
		return
	}

	cookie := http.Cookie{
		Name:     "email",
		Value:    email,
		Expires:  time.Now().UTC().AddDate(0, 6, 0),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("Activated"))
}

func (e *Emails) Activate(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	token := vars["token"]

	email, err := data.ActivateEmail(token)
	if err != nil {
		e.l.Printf("[ERROR] activating email: %v", err)
		http.Error(w, "Failed activating email", http.StatusBadRequest)
		return
	}

	resp, err := http.Post("https://localhost:8083/verifyemail", "text/http", bytes.NewBuffer([]byte(email)))
	if err != nil {
		e.l.Printf("[ERROR] activating email: %v", err)
		http.Error(w, "Failed activating email", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		e.l.Println("[ERROR] activating email")
		http.Error(w, "Failed activating email", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Email activated"))
}