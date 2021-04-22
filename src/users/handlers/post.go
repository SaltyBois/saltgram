package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"saltgram/users/data"
	"time"

	"github.com/dgrijalva/jwt-go"
)

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

		go data.SendActivation(user.Email)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Activation email sent"))
		// if err != nil {
		// 	u.l.Printf("[ERROR] sending activation: %v\n", err)
		// 	http.Error(w, "Failed to send activation for user", http.StatusInternalServerError)
		// 	return
		// }
	}
}
