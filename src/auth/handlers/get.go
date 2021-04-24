package handlers

import (
	"net/http"
	"os"
	"saltgram/auth/data"
	"strconv"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
)

func (a *Auth) CheckPermissions(e *casbin.Enforcer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(Jovan): Get from user srvc

	}
}

func (a *Auth) Refresh(db *data.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		refreshToken, err := data.GetRefreshToken(db, claims.Username)

		if err != nil {
			a.l.Println("[ERROR] can't find refresh token")
			http.Error(w, "Can't find refresh token", http.StatusBadRequest)
			return
		}

		rt := data.Refresh{
			Username: claims.Username,
			Token:    refreshToken,
		}
		if err := rt.Verify(db); err != nil {
			a.l.Println("[ERROR] refresh token no longer valid")
			http.Error(w, "refresh token no longer valid", http.StatusBadRequest)
			return
		}

		//
		jws, err := getJWS(r)
		if err != nil {
			a.l.Println("[ERROR] JWS not found")
			http.Error(w, "Missing JWS", http.StatusBadRequest)
			return
		}

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
}
