package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"saltgram/auth/data"
	"time"

	"github.com/dgrijalva/jwt-go"
)



func (a *Auth) AddRefreshToken(db *data.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Context().Value(KeyRefreshToken{}).(data.Refresh)
		err := data.AddRefreshToken(db, token.Username, token.Token)
		if err != nil {
			a.l.Printf("[ERROR] adding refresh token: %v\n", err)
			http.Error(w, "Error adding refresh token", http.StatusBadRequest)
			return
		}
	}
}

func (a *Auth) GetJWT(db *data.DBConn) func(http.ResponseWriter, *http.Request){
	return func (w http.ResponseWriter, r *http.Request) {
		user := data.Login{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			a.l.Printf("[ERROR] deserializing user: %v", err)
			http.Error(w, "Failed to deserialize user", http.StatusBadRequest)
			return
		}

		// NOTE(Jovan): HS256 is considered safe enough
		claims := AccessClaims{
			Username: user.Username,
			Password: user.Password,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().UTC().Add(time.Second * 5).Unix(),
				Issuer:    "SaltGram",
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		jws, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
		if err != nil {
			a.l.Printf("[ERROR] failed signing JWT: %v", err)
			http.Error(w, "Failed signing JWT: "+err.Error(), http.StatusInternalServerError)
			return
		}

		refreshToken, err := data.GetRefreshToken(db, user.Username)
		if err != nil {
			a.l.Printf("[ERROR] failed getting refresh token: %v", err)
			http.Error(w, "Failed to get refresh token", http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     "refresh",
			Value:    refreshToken,
			Expires:  time.Now().UTC().AddDate(0, 6, 0),
			HttpOnly: true,
		}

		http.SetCookie(w, &cookie)

		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte(jws))
	}
}