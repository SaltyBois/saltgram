package handlers

import (
	"io/ioutil"
	"net/http"
	"saltgram/data"

	"github.com/go-playground/validator"
)

type ChangeRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

func (rr *ChangeRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(rr)
}

func (e *Emails) ChangePassword(w http.ResponseWriter, r *http.Request) {
	cr := ChangeRequest{}
	err := data.FromJSON(&cr, r.Body)
	if err != nil {
		e.l.Printf("[ERROR] deserializing ChangeRequest: %v\n", err)
		http.Error(w, "Failed to parse request", http.StatusBadRequest)
		return
	}

	err = cr.Validate()
	if err != nil {
		e.l.Printf("[ERROR] ChangeRequest not valid: %v\n", err)
		http.Error(w, "Bad change request", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("email")
	if err != nil {
		e.l.Printf("[ERROR] getting cookie: %v", err)
		http.Error(w, "No change request cookie", http.StatusBadRequest)
		return
	}
	email := cookie.Value
	err = data.ChangePassword(email, cr.OldPassword, cr.NewPassword)
	if err != nil {
		e.l.Printf("[ERROR] changing password: %v\n", err)
		http.Error(w, "Failed to change password", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
}

func (e *Emails) RequestReset(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e.l.Printf("[ERROR] getting email: %v\n", err)
		http.Error(w, "No email", http.StatusBadRequest)
		return
	}
	email := string(body)
	err = data.SendPasswordReset(email)
	if err != nil {
		e.l.Printf("[ERROR] sending email request: %v\n", err)
	}
	// NOTE(Jovan): Always return 200 OK as per OWASP guidelines
}

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
