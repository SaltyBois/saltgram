package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"saltgram/users/data"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
)

type ChangeRequest struct {
	Email string `json:"email" validate: "required"`
	OldPlainPassword string `json:"oldPlainPassword" validate:"required`
	NewPlainPassword string `json:"newPlainPassword" validate:"required"`
}

func (cr *ChangeRequest) Validate() error {
	valid := validator.New()
	return valid.Struct(cr)
}

type KeyChangeRequest struct{}

func (u *Users) VerifyEmail(db *data.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u.l.Print("Activating email")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			u.l.Printf("[ERROR] getting email: %v\n", err)
			http.Error(w, "No email", http.StatusBadRequest)
			return
		}
		email := string(body)

		err = data.VerifyEmail(db, email)
		if err != nil {
			u.l.Printf("[ERROR] verifying email: %v\n", err)
			http.Error(w, "Failed to verify email", http.StatusBadRequest)
			return
		}

		w.Write([]byte("Verified email"))
	}
}

func (u *Users) ChangePassword(db *data.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u.l.Println("Changing password")

		changeRequest := r.Context().Value(KeyChangeRequest{}).(ChangeRequest)

		err := data.ChangePassword(db, changeRequest.Email, changeRequest.OldPlainPassword, changeRequest.NewPlainPassword)
		if err != nil {
			u.l.Printf("[ERROR] attepmting to change password: %v\n")
			http.Error(w, "Failed to change password", http.StatusBadRequest)
			return
		}
		w.Write([]byte("Password changed"))
	}
}

func (u *Users) Register(db *data.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u.l.Println("Handling POST Users")

		user := r.Context().Value(KeyUser{}).(data.User)
		refreshClaims := RefreshClaims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				// TODO(Jovan): Make programmatic?
				ExpiresAt: time.Now().UTC().AddDate(0, 6, 0).Unix(),
				Issuer:    "SaltGram",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
		jws, err := token.SignedString([]byte(os.Getenv("REF_SECRET_KEY")))

		if err != nil {
			u.l.Println("[ERROR] signing refresh token")
			http.Error(w, "Failed signing refresh token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = db.AddUser(&user)
		if err != nil {
			u.l.Printf("[ERROR] adding user: %v\n", err)
			http.Error(w, "Failed to register user", http.StatusBadRequest)
			return
		}

		// NOTE(Jovan): Saving refresh token
		// data.AddRefreshToken(user.Username, jws)
		values := map[string]string{"username": user.Username, "token": jws}
		jsonData, err := json.Marshal(values)
		if err != nil {
			u.l.Printf("[ERROR] marshalling refresh token request: %v\n", err)
			http.Error(w, "Error marshalling refresh token request", http.StatusInternalServerError)
			return
		}

		_, err = http.Post("https://localhost:8082/refresh", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			u.l.Printf("[ERROR] POST refresh token: %v\n", err)
			http.Error(w, "Error in POST refresh token", http.StatusInternalServerError)
			return
		}
		// TODO(Jovan): Check response???

		// go data.SendActivation(user.Email)
		go func() {
			_, err = http.Post("https://localhost:8084/activate", "text/html", bytes.NewBuffer([]byte(user.Email)))
			if err != nil {
				u.l.Printf("[ERROR] sending change password request: %v\n", err)
			}
		}()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Activation email sent"))
		// if err != nil {
		// 	u.l.Printf("[ERROR] sending activation: %v\n", err)
		// 	http.Error(w, "Failed to send activation for user", http.StatusInternalServerError)
		// 	return
		// }
	}
}
