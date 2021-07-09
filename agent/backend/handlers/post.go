package handlers

import (
	"agent/data"
	"net/http"
)

func (a *Agent) SignIn(w http.ResponseWriter, r *http.Request) {
	dto := SignInDTO{}
	FromJSON(&dto, r.Body)

	jws, err := a.db.Login(dto.Username, dto.Password)
	if err != nil {
		a.l.Errorf("failed to login: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(jws))
}

func (a *Agent) Signup(w http.ResponseWriter, r *http.Request) {
	dto := UserDTO{}
	FromJSON(&dto, r.Body)

	err := a.db.AddUser(&data.User{
		Username: dto.Username,
		Password: dto.Password,
		Email:    dto.Email,
		Agent:    dto.Agent,
		Token:    dto.Token,
	})
	if err != nil {
		a.l.Errorf("failed to register user: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.Write([]byte("OK"))
}
