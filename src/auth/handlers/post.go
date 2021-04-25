package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"saltgram/auth/data"
	saltdata "saltgram/data"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (l *Login) Login(w http.ResponseWriter, r *http.Request) {
	l.l.Println("Handling POST reCaptcha")
	login := r.Context().Value(KeyLogin{}).(saltdata.Login)
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

	// if !data.IsEmailVerified(login.Username) {
	// 	http.Error(w, "Email not activated", http.StatusForbidden)
	// 	return
	// }

	resp, err := http.Get(fmt.Sprintf("https://localhost:%s/verifyemail/%s", os.Getenv("SALT_USERS_PORT"), login.Username))
	if err != nil {
		l.l.Printf("[ERROR] checking user email: %v\n", err)
		http.Error(w, "Error checking user email", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		l.l.Printf("[ERROR] Email not verified")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// NOTE(Jovan): Returning reCaptcha score for testing purposes
	// err = data.ToJSON(fmt.Sprintf("reCAPTCHA score: %f", score), w)
	// if err != nil {
	// 	http.Error(w, "Error serializing score somehow: "+err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// hashedPass, err := data.VerifyPassword(login.Username, login.Password)
	values := map[string]string{"username": login.Username, "password": login.Password}
	jsonData, err := json.Marshal(values)
	if err != nil {
		l.l.Printf("[ERROR] marshalling password check request: %v\n", err)
		http.Error(w, "Error marshalling password check request", http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://localhost:%s/password", os.Getenv("SALT_USERS_PORT")), bytes.NewBuffer(jsonData))
	if err != nil {
		l.l.Printf("[ERROR] creating PUT request: %v\n", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	client := http.Client{}
	resp, err = client.Do(req)

	if err != nil {
		l.l.Printf("[ERROR] PUT method: %v\n", err)
		http.Error(w, "Failed to do PUT method", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		l.l.Println("[ERROR] invalid password")
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.l.Printf("[ERROR] reading PUT response: %v\n", err)
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}
	hashedPass := string(body)

	w.Header().Add("Content-Type", "application/json")
	u := saltdata.Login{Username: login.Username, Password: hashedPass}
	err = saltdata.ToJSON(u, w)
	if err != nil {
		l.l.Printf("[ERROR] serializing login: %v", err)
		http.Error(w, "Failed to serialize login", http.StatusInternalServerError)
		return
	}
}

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

func (a *Auth) GetJWT(db *data.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := saltdata.Login{}
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
