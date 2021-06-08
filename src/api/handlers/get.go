package handlers

import (
	"context"
	"io"
	"net/http"
	"os"
	saltdata "saltgram/data"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/users/prusers"

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

	jws, err := getUserJWS(r)
	if err != nil {
		a.l.Println("[ERROR] JWS not found")
		http.Error(w, "Missing JWS", http.StatusBadRequest)
		return
	}

	res, err := a.ac.Refresh(context.Background(), &prauth.RefreshRequest{OldJWS: jws, Refresh: cookie.Value})
	if err != nil {
		a.l.Printf("[ERROR] getting refresh token: %v\n", err)
		http.Error(w, "Failed to get refresh token", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(res.NewJWS))
}

func (u *Users) GetByJWS(w http.ResponseWriter, r *http.Request) {
	jws, err := getUserJWS(r)
	if err != nil {
		u.l.Println("[ERROR] JWS not found")
		http.Error(w, "JWS not found", http.StatusBadRequest)
		return
	}

	token, err := jwt.ParseWithClaims(
		jws,
		&saltdata.AccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		},
	)

	if err != nil {
		u.l.Printf("[ERROR] parsing claims: %v", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		u.l.Println("[ERROR] unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	user, err := u.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: claims.Username})
	if err != nil {
		u.l.Println("[ERROR] fetching user", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.HashedPassword != claims.Password {
		u.l.Println("[ERROR] passwords do not match")
		http.Error(w, "JWT password doesn't match user's password", http.StatusUnauthorized)
		return
	}

	err = saltdata.ToJSON(user, w)
	if err != nil {
		u.l.Println("[ERROR] serializing user ", err)
		http.Error(w, "Error serializing user", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
}

func (u *Users) GetProfile(w http.ResponseWriter, r *http.Request) {

	jws, err := getUserJWS(r)
	if err != nil {
		u.l.Println("[ERROR] JWS not found")
		http.Error(w, "JWS not found", http.StatusBadRequest)
		return
	}

	token, err := jwt.ParseWithClaims(
		jws,
		&saltdata.AccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		},
	)

	if err != nil {
		u.l.Printf("[ERROR] parsing claims: %v", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		u.l.Println("[ERROR] unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	user := claims.Username

	vars := mux.Vars(r)
	profileUsername, er := vars["username"]
	if !er {
		u.l.Println("[ERROR] parsing URL, no username in URL")
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	profile, err := u.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: user, Username: profileUsername})
	if err != nil {
		u.l.Println("[ERROR] fetching profile")
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	saltdata.ToJSON(profile, w)

}

func (u *Users) GetFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, er := vars["username"]
	if !er {
		u.l.Println("[ERROR] parsing URL, no username in URL")
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	stream, err := u.uc.GetFollowers(context.Background(), &prusers.FollowerRequest{Username: username})
	if err != nil {
		u.l.Println("[ERROR] fetching followers")
		http.Error(w, "Followers fetching error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("{"))
	for {
		profile, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			u.l.Println("[ERROR] fetching followers")
			http.Error(w, "Error couldn't fetch followers", http.StatusInternalServerError)
			return
		}
		saltdata.ToJSON(profile, w)
	}
	w.Write([]byte("}"))
}

func (u *Users) GetFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, er := vars["username"]
	if !er {
		u.l.Println("[ERROR] parsing URL, no username in URL")
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	stream, err := u.uc.GetFollowers(context.Background(), &prusers.FollowerRequest{Username: username})
	if err != nil {
		u.l.Println("[ERROR] fetching followers")
		http.Error(w, "Followers fetching error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("{"))
	for {
		profile, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			u.l.Println("[ERROR] fetching followers")
			http.Error(w, "Error couldn't fetch followers", http.StatusInternalServerError)
			return
		}
		saltdata.ToJSON(profile, w)
	}
	w.Write([]byte("}"))
}
