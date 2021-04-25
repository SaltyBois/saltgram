package handlers

import (
	"context"
	"net/http"
	saltdata "saltgram/data"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/users/prusers"
)

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
