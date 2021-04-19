package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"saltgram/data"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
)

type ChangeRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

func (rr *ChangeRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(rr)
}

func (e *Emails) ChangePassword(w http.ResponseWriter, r *http.Request) {
	cr := ChangeRequest{}
	err := data.FromJSON(&cr, r.Body)
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
	email := cookie.Value
	err = data.ChangePassword(email, cr.OldPassword, cr.NewPassword)
	if err != nil {
		e.l.Printf("[ERROR] changing password: %v\n", err)
		http.Error(w, "Failed to change password", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
}

func (e *Emails) RequestReset(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e.l.Printf("[ERROR] getting email: %v\n", err)
		http.Error(w, "No email", http.StatusBadRequest)
		return
	}
	email := string(body)
	err = data.SendPasswordReset(email)
	if err != nil {
		e.l.Printf("[ERROR] sending email request: %v\n", err)
	}
	// NOTE(Jovan): Always return 200 OK as per OWASP guidelines
}

func (a *Auth) GetJWT(w http.ResponseWriter, r *http.Request) {
	user := data.Login{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		a.l.Printf("[ERROR] deserializing user: %v", err)
		http.Error(w, "Failed to deserialize user", http.StatusBadRequest)
		return
	}

	// TODO(Jovan): Pull out into const
	// timeout, _ := strconv.Atoi(os.Getenv("TOKEN_TIMEOUT_MINUTES"))

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

	refreshToken, err := data.GetRefreshToken(user.Username)
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

	if !data.IsEmailVerified(login.Username) {
		http.Error(w, "Email not activated", http.StatusForbidden)
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

	w.Header().Add("Content-Type", "application/json")
	u := data.Login{Username: login.Username, Password: hashedPass}
	err = data.ToJSON(u, w)
	if err != nil {
		l.l.Printf("[ERROR] serializing login: %v", err)
		http.Error(w, "Failed to serialize login", http.StatusInternalServerError)
		return
	}
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
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
	data.AddUser(&user)
	data.AddRefreshToken(user.Username, jws)
	go data.SendActivation(user.Email)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Activation email sent"))
	// if err != nil {
	// 	u.l.Printf("[ERROR] sending activation: %v\n", err)
	// 	http.Error(w, "Failed to send activation for user", http.StatusInternalServerError)
	// 	return
	// }

}
