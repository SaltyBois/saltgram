package handlers

import (
	"fmt"
	"net/http"
	"saltgram/data"
)

func (re *Login) Login(w http.ResponseWriter, r *http.Request) {
	re.l.Println("Handling POST reCaptcha")
	reCaptcha := r.Context().Value(KeyLogin{}).(data.Login)
	score, err := data.VerifyCaptcha(&reCaptcha)
	if err != nil {
		re.l.Println("[ERROR] verifying reCaptcha")
		http.Error(w, "Failed verifying reCaptcha: "+err.Error(), http.StatusBadRequest)
		return
	}

	if score < 0.5 {
		http.Error(w, "Low reCaptcha score", http.StatusBadRequest)
		return
	}

	err = data.ToJSON(fmt.Sprintf("reCAPTCHA score: %f", score), w)
	if err != nil {
		http.Error(w, "Error serializing score somehow: "+err.Error(), http.StatusBadRequest)
		return
	}
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	u.l.Println("Logging in...")
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	u.l.Println("Handling POST Users")

	user := r.Context().Value(KeyUser{}).(data.User)
	data.AddUser(&user)
}
