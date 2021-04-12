package handlers

import (
	"net/http"
	"os"
	"saltgram/data"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func (e *Emails) Activate(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	token := vars["token"]

	err := data.ActivateEmail(token)
	if err != nil {
		e.l.Printf("[ERROR] activating email: %v", err)
		http.Error(w, "Failed activating email: "+err.Error(), http.StatusInternalServerError)
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

func (u *Users) GetByJWS(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	jws := getUserJWS(r)

	token, err := jwt.ParseWithClaims(
		jws,
		&UserClaims{},
		func(t *jwt.Token)(interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		},
	)

	if err != nil {
		u.l.Printf("[ERROR] parsing claims: %v", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*UserClaims)

	if !ok {
		u.l.Println("[ERROR] unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	if err := claims.CheckDate(); err != nil {
		u.l.Println("[ERROR] jwt expired")
		http.Error(w, "JWT expired", http.StatusUnauthorized)
		return
	}

	user, err := data.GetUserByUsername(claims.Username)
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

	if user.HashedPassword != claims.Password {
		u.l.Println("[ERROR] passwords do not match")
		http.Error(w, "JWT password doesn't match user's password", http.StatusUnauthorized)
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
