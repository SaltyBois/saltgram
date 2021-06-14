package handlers

import (
	"context"
	"io/ioutil"
	"net/http"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/email/premail"
	"time"

	"github.com/gorilla/mux"
)

func (a *Auth) CheckPermissions(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.l.Printf("[ERROR] reading body: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	route := string(body)
	jws, _ := getUserJWS(r)
	_, err = a.ac.CheckPermissions(context.Background(),
		&prauth.PermissionRequest{
			Jws:    jws,
			Path:   route,
			Method: r.Method,
		})
	if err != nil {
		a.l.Printf("[ERROR] authenticating: %v\n", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	w.Write([]byte("200 - OK"))
}

func (e *Email) ConfirmReset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	res, err := e.ec.ConfirmReset(context.Background(), &premail.ConfirmRequest{Token: token})
	if err != nil {
		e.l.Errorf("failure confirming reset: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	cookie := http.Cookie{
		Name:     "emailforreset",
		Value:    res.Email,
		Expires:  time.Now().UTC().AddDate(0, 6, 0),
		HttpOnly: true,
		Secure:   true,
		Path:     "/users",
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("200 - OK"))
}

func (e *Email) Activate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	_, err := e.ec.Activate(context.Background(), &premail.ActivateRequest{Token: token})
	if err != nil {
		e.l.Errorf("failure activating email: %v\n", err)
		http.Error(w, "Failed to activate email", http.StatusBadRequest)
		return
	}

	w.Write([]byte("200 - OK"))
}
