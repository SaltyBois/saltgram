package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"saltgram/email/data"

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

func (e *Emails) SendActivation(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e.l.Printf("[ERROR] getting email: %v\n", err)
		http.Error(w, "No email", http.StatusBadRequest)
		return
	}
	email := string(body)

	err = data.SendActivation(email)
	if err != nil {
		e.l.Printf("[ERROR] sending email activation: %v\n", err)
		http.Error(w, "Failed to send email activation", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Activation sent"))
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
	// err = data.ChangePassword(email, cr.OldPassword, cr.NewPassword)
	values := map[string]string{"email": email, "oldPlainPassword": cr.OldPassword, "newPlainPassword": cr.NewPassword}
	jsonData, err := json.Marshal(values)
	if err != nil {
		e.l.Printf("[ERROR] marshalling change password request: %v\n", err)
		http.Error(w, "Error marshalling change password request", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("https://localhost:8083/password", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		e.l.Printf("[ERROR] POST change password request: %v\n", err)
		http.Error(w, "Error in POST change password request", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		e.l.Println("[ERROR] failed to change password")
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
