package handlers

import (
	"context"
	"io"
	"net/http"
	"os"
	saltdata "saltgram/data"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/content/prcontent"
	"saltgram/protos/users/prusers"

	"github.com/dgrijalva/jwt-go"
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

func (s *Content) GetSharedMedia(w http.ResponseWriter, r *http.Request) {

	jws, err := getUserJWS(r)
	if err != nil {
		s.l.Println("[ERROR] JWS not found")
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
		s.l.Printf("[ERROR] parsing claims: %v", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		s.l.Println("[ERROR] unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	user, err := s.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: claims.Username})
	if err != nil {
		s.l.Println("[ERROR] fetching user", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	stream, err := s.cc.GetSharedMedia(context.Background(), &prcontent.SharedMediaRequest{UserId: user.Id})
	if err != nil {
		s.l.Println("[ERROR] fetching shared medias")
		http.Error(w, "Followers shared media error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("{"))
	for {
		sharedMedia, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.l.Println("[ERROR] fetching sharedMedias")
			http.Error(w, "Error couldn't fetch sharedMedias", http.StatusInternalServerError)
			return
		}
		saltdata.ToJSON(sharedMedia, w)
	}
	w.Write([]byte("}"))
}
