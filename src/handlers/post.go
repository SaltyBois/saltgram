package handlers

import (
	"net/http"
	"os"
	"saltgram/data"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (l *Login) Login(w http.ResponseWriter, r *http.Request) {
	l.l.Println("Handling POST reCaptcha")
	login := r.Context().Value(KeyLogin{}).(data.Login)
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

	// NOTE(Jovan): Returning reCaptcha score for testing purposes
	// err = data.ToJSON(fmt.Sprintf("reCAPTCHA score: %f", score), w)
	// if err != nil {
	// 	http.Error(w, "Error serializing score somehow: "+err.Error(), http.StatusBadRequest)
	// 	return
	// }

	hashedPass, err := data.VerifyPassword(login.Username, login.Password)
	if err != nil {
		l.l.Println("[ERROR] invalid password")
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// NOTE(Jovan): HS256 is considered safe enough
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": login.Username,
		"password": hashedPass,
		"nbf":      time.Now().UTC().AddDate(0, 0, 1).String(), // TODO(Jovan): Change to appropriate time duration
	})

	jws, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil{
		l.l.Printf("[ERROR] failed signing JWT: %v", err)
		http.Error(w, "Failed signing JWT: " + err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(jws))

}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	u.l.Println("Handling POST Users")

	user := r.Context().Value(KeyUser{}).(data.User)
	data.AddUser(&user)
	data.SendActivation(user.Email)
}
