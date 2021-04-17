package handlers

import (
	"net/http"
	"os"
	"saltgram/data"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func (a *Auth) Refresh(w http.ResponseWriter, r *http.Request) {
	a.l.Println("Handling REFRESH")
	cookie, err := r.Cookie("refresh")
	if err != nil {
		a.l.Printf("[ERROR] getting cookie: %v", err)
		http.Error(w, "No refresh cookie", http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value
	token, err := jwt.ParseWithClaims(
		tokenString,
		&RefreshClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("REF_SECRET_KEY")), nil
		},
	)

	if err != nil {
		a.l.Printf("[ERROR] parsing refresh claims: %v", err)
		http.Error(w, "Invalid refresh cookie", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*RefreshClaims)

	if !ok {
		a.l.Println("[ERROR] unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	refreshToken, err := data.GetRefreshToken(claims.Username)

	if err != nil {
		a.l.Println("[ERROR] can't find refresh token")
		http.Error(w, "Can't find refresh token", http.StatusBadRequest)
		return
	}

	rt := data.Refresh{
		Username: claims.Username,
		Token:    refreshToken,
	}
	if err := rt.Validate(); err != nil {
		a.l.Println("[ERROR] refresh token no longer valid")
		http.Error(w, "refresh token no longer valid", http.StatusBadRequest)
		return
	}

	//
	jws := getUserJWS(r)

	// NOTE(Jovan): Not validating 'cause it is invalid
	jwtOld, _ := jwt.ParseWithClaims(
		jws,
		&AccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		},
	)

	jwsClaims, ok := jwtOld.Claims.(*AccessClaims)

	if !ok {
		a.l.Println("[ERROR] unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	// TODO(Jovan): Pull out into const
	timeout, _ := strconv.Atoi(os.Getenv("TOKEN_TIMEOUT_MINUTES"))
	jwsClaims.StandardClaims.ExpiresAt = time.Now().UTC().Add(time.Minute * time.Duration(timeout)).Unix()
	jwtNew := jwt.NewWithClaims(jwt.SigningMethodHS256, jwsClaims)

	jwsNew, err := jwtNew.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		a.l.Printf("[ERROR] failed signing JWT: %v", err)
		http.Error(w, "Failed signing JWT: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(jwsNew))
}

func (e *Emails) ConfirmReset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	email, err := data.ConfirmPasswordReset(token)
	if err != nil {
		e.l.Printf("[ERROR] confirming password reset: %v\n", err)
		http.Error(w, "Failed to confirm password reset", http.StatusBadRequest)
		return
	}

	cookie := http.Cookie{
		Name:     "email",
		Value:    email,
		Expires:  time.Now().UTC().AddDate(0, 6, 0),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("200 - OK"))
}

// TODO(Jovan): REMOVE!
func (e *Emails) GetAllResets(w http.ResponseWriter, r *http.Request) {
	data.ToJSON(data.GetAllResets(), w)
}

func (e *Emails) Activate(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	token := vars["token"]

	err := data.ActivateEmail(token)
	if err != nil {
		e.l.Printf("[ERROR] activating email: %v", err)
		http.Error(w, "Failed activating email: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Email activated"))
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
		&AccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		},
	)

	if err != nil {
		u.l.Printf("[ERROR] parsing claims: %v", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*AccessClaims)

	if !ok {
		u.l.Println("[ERROR] unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
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
