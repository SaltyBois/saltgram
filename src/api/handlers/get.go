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
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func (a *Auth) Authenticate2FA(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	_, err := a.ac.Authenticate2FA(context.Background(), &prauth.Auth2FARequest{Token: token})
	if err != nil {
		a.l.Errorf("failed to authenticate 2fa token: %v, error: %v\n", token, err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.Write([]byte("Authenticated!"))
}

func (a *Auth) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh")
	if err != nil {
		a.l.Printf("getting cookie: %v", err)
		http.Error(w, "No refresh cookie", http.StatusBadRequest)
		return
	}

	jws, err := getUserJWS(r)
	if err != nil {
		a.l.Errorf("failed to get user jws: %v\n", err)
		http.Error(w, "Missing JWS", http.StatusBadRequest)
		return
	}

	a.l.Printf("Refreshing token: %v\n", jws)

	res, err := a.ac.Refresh(context.Background(), &prauth.RefreshRequest{OldJWS: jws, Refresh: cookie.Value})
	if err != nil {
		a.l.Errorf("getting refresh token: %v\n", err)
		http.Error(w, "Failed to get refresh token", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(res.NewJWS))
}

func (u *Users) GetByJWS(w http.ResponseWriter, r *http.Request) {
	jws, err := getUserJWS(r)
	if err != nil {
		u.l.Errorf("failed to get user's jws: %v\n", err)
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
		u.l.Errorf("failed to parse claims: %v", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		u.l.Errorf("unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	user, err := u.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: claims.Username})
	if err != nil {
		u.l.Errorf("failed to fetch user: %v\n", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.HashedPassword != claims.Password {
		u.l.Println("passwords do not match")
		http.Error(w, "JWT password doesn't match user's password", http.StatusUnauthorized)
		return
	}

	err = saltdata.ToJSON(user, w)
	if err != nil {
		u.l.Errorf("serializing user ", err)
		http.Error(w, "Error serializing user", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
}

func (s *Content) GetSharedMedia(w http.ResponseWriter, r *http.Request) {

	jws, err := getUserJWS(r)
	if err != nil {
		s.l.Println("getting jws: %v\n", err)
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
		s.l.Errorf("parsing claims: %v", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		s.l.Errorf("unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	user, err := s.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: claims.Username})
	if err != nil {
		s.l.Errorf("failed to fetch user: %v\n", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	stream, err := s.cc.GetSharedMedia(context.Background(), &prcontent.SharedMediaRequest{UserId: user.Id})
	if err != nil {
		s.l.Errorf("failed to fetch media %v\n", err)
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
			s.l.Errorf("failed to fetch sharedMedias: %v\n", err)
			http.Error(w, "Error couldn't fetch sharedMedias", http.StatusInternalServerError)
			return
		}
		saltdata.ToJSON(sharedMedia, w)
	}
	w.Write([]byte("}"))
}

func (s *Content) GetSharedMediaByUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userId := vars["id"]
	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		s.l.Println("converting id")
		return
	}

	stream, err := s.cc.GetSharedMedia(context.Background(), &prcontent.SharedMediaRequest{UserId: id})
	if err != nil {
		s.l.Errorf("failed fetching shared medias %v\n", err)
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
			s.l.Errorf("failed to fetch sharedMedias: %v\n", err)
			http.Error(w, "Error couldn't fetch sharedMedias", http.StatusInternalServerError)
			return
		}
		saltdata.ToJSON(sharedMedia, w)
	}
	w.Write([]byte("}"))
}

func (c *Content) GetProfilePictureByUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userId := vars["id"]
	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		c.l.Println("[ERROR] converting id")
		return
	}

	profilePicture, err := c.cc.GetProfilePicture(context.Background(), &prcontent.GetProfilePictureRequest{UserId: id})
	if err != nil {
		c.l.Println("[ERROR] fetching profile picture", err)
		http.Error(w, "pp not found", http.StatusNotFound)
		return
	}

	err = saltdata.ToJSON(profilePicture, w)
	if err != nil {
		c.l.Println("[ERROR] serializing pp ", err)
		http.Error(w, "Error serializing pp", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")

}
