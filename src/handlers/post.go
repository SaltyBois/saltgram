package handlers

import (
	"net/http"
	"saltgram/data"
)

func (l *Login) Login(w http.ResponseWriter, r *http.Request) {
	l.l.Println("Handling POST reCaptcha")
	login := r.Context().Value(KeyLogin{}).(data.Login)
	score, err := login.ReCaptcha.Verify()
	if err != nil {
		l.l.Println("[ERROR] verifying reCaptcha")
		http.Error(w, "Failed verifying reCaptcha: "+err.Error(), http.StatusBadRequest)
		return
	}

	if score < 0.5 {
		http.Error(w, "Low reCaptcha score", http.StatusBadRequest)
		return
	}

	if !data.IsEmailVerified(login.Username) {
		http.Error(w, "Email not activated", http.StatusForbidden)
		return
	}

	// NOTE(Jovan): Returning reCaptcha score for testing purposes
	// err = data.ToJSON(fmt.Sprintf("reCAPTCHA score: %f", score), w)
	// if err != nil {
	// 	http.Error(w, "Error serializing score somehow: "+err.Error(), http.StatusBadRequest)
	// 	return
	// }

	hashedPass, err := data.VerifyPassword(login.Username, login.Password)
	if err != nil {
		l.l.Println("[ERROR] invalid password")
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	u := data.Login{Username: login.Username, Password: hashedPass}
	err = data.ToJSON(u, w)
	if err != nil {
		l.l.Printf("[ERROR] serializing login: %v", err)
		http.Error(w, "Failed to serialize login", http.StatusInternalServerError)
		return
	}
}
