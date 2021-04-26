package handlers

import (
	"context"
	"io/ioutil"
	"net/http"
	saltdata "saltgram/data"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"

	"github.com/go-playground/validator"
)

type ChangeRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

func (cr *ChangeRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(cr)
}

func (e *Email) ChangePassword(w http.ResponseWriter, r *http.Request) {
	cr := ChangeRequest{}

	err := saltdata.FromJSON(&cr, r.Body)
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
	_, err = e.uc.ChangePassword(context.Background(), &prusers.ChangeRequest{
		Email: cookie.Value,
		OldPlainPassword: cr.OldPassword,
		NewPlainPassword: cr.NewPassword,
	})
	if err != nil {
		e.l.Printf("[ERROR] POST change password request: %v\n", err)
		http.Error(w, "Error in POST change password request", http.StatusBadRequest)
		return
	}
	w.Write([]byte("200 - OK"))
}

func (e *Email) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e.l.Printf("[ERROR] getting email: %v\n", err)
		http.Error(w, "No email", http.StatusBadRequest)
		return
	}
	email := string(body)
	_, err = e.ec.RequestReset(context.Background(), &premail.ResetRequest{Email: email})
	if err != nil {
		e.l.Printf("[ERROR] sending email request: %v\n", err)
	}
	// NOTE(Jovan): Always send 200 OK as per OWASP
	w.Write([]byte("200 - OK"))
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	dto := saltdata.UserDTO{}
	err := saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		u.l.Printf("[ERROR] deserializing user data: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = dto.Validate()
	if err != nil {
		u.l.Printf("[ERROR] validating user data: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = u.uc.Register(context.Background(), &prusers.RegisterRequest{
		Username: dto.Username,
		FullName: dto.FullName,
		Email:    dto.Email,
		Password: dto.Password,
		ReCaptcha: &prusers.UserReCaptcha{
			Token:  dto.ReCaptcha.Token,
			Action: dto.ReCaptcha.Action,
		},
	})

	if err != nil {
		u.l.Printf("[ERROR] registering user: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Activation email sent"))
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	login := saltdata.Login{}
	err := saltdata.FromJSON(&login, r.Body)
	if err != nil {
		a.l.Printf("[ERROR] deserializing body: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = login.Validate()
	if err != nil {
		a.l.Printf("[ERROR] validating: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	res, err := a.authClient.Login(context.Background(), &prauth.LoginRequest{
		Username: login.Username,
		Password: login.Password,
		ReCaptcha: &prauth.ReCaptcha{
			Action: login.ReCaptcha.Action,
			Token:  login.ReCaptcha.Token,
		},
	})

	if err != nil {
		a.l.Printf("[ERROR] calling login: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	saltdata.ToJSON(res, w)
}
