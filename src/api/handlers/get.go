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

	u.l.Println("[INFO] username: %v\n", claims.Username)
	u.l.Println("[INFO] password: %v\n", claims.Password)

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

	response := saltdata.UserDTO{
		Id:       strconv.FormatUint(user.Id, 10),
		Email:    user.Email,
		FullName: user.FullName,
		Username: user.Username,
	}

	err = saltdata.ToJSON(response, w)
	if err != nil {
		u.l.Errorf("serializing user ", err)
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

	saltdata.ToProtoJSON(profile, w)

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
	var profiles []*prusers.ProfileFollower
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
		profiles = append(profiles, profile)
	}
	saltdata.ToJSON(profiles, w)
}

func (u *Users) GetFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, er := vars["username"]
	if !er {
		u.l.Println("[ERROR] parsing URL, no username in URL")
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	stream, err := u.uc.GerFollowing(context.Background(), &prusers.FollowerRequest{Username: username})
	if err != nil {
		u.l.Println("[ERROR] fetching following", err)
		http.Error(w, "Following fetching error", http.StatusInternalServerError)
		return
	}
	var profiles []*prusers.ProfileFollower
	for {
		profile, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			u.l.Println("[ERROR] fetching following", err)
			http.Error(w, "Error couldn't fetch following", http.StatusInternalServerError)
			return
		}
		profiles = append(profiles, profile)
	}
	saltdata.ToJSON(profiles, w)
}

func (u *Users) GetFollowingRequest(w http.ResponseWriter, r *http.Request) {
	jws, err := getUserJWS(r)
	if err != nil {
		u.l.Println("getting jws: %v\n", err)
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
		u.l.Errorf("parsing claims: %v", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		u.l.Errorf("unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	stream, err := u.uc.GetFollowRequests(context.Background(), &prusers.Profile{Username: claims.Username})
	if err != nil {
		u.l.Println("[ERROR] fetching follower request")
		http.Error(w, "Follower request fetching error", http.StatusInternalServerError)
		return
	}

	var profiles []*prusers.FollowingRequest
	for {
		profile, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			u.l.Println("[ERROR] fetching following", err)
			http.Error(w, "Error couldn't fetch following", http.StatusInternalServerError)
			return
		}
		profiles = append(profiles, profile)
	}
	saltdata.ToJSON(profiles, w)

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

func (u *Users) SearchUsers(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	queryUsername := vars["username"]

	queryResults, err := u.uc.GetSearchedUsers(context.Background(), &prusers.SearchRequest{Query: queryUsername})

	if err != nil {
		u.l.Println("[ERROR] Searching users failed")
		http.Error(w, "Error Searching users: ", http.StatusInternalServerError)
		return
	}

	finalResult := []*prusers.SearchedUser{}

	const MAX_NUMBER_OF_RESULTS = 20

	for i := 0; i < len(queryResults.SearchedUser); i++ {
		su := queryResults.SearchedUser[i]
		if su.Username == claims.Username {
			continue
		}
		if i == MAX_NUMBER_OF_RESULTS {
			break
		}
		finalResult = append(finalResult, &prusers.SearchedUser{
			Username:              su.Username,
			ProfilePictureAddress: su.ProfilePictureAddress})

	}

	err = saltdata.ToJSON(&finalResult, w)

	if err != nil {
		u.l.Println("[ERROR] Searching users failed")
		http.Error(w, "Error Searching users: ", http.StatusInternalServerError)
		return
	}

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

	w.Write([]byte(profilePicture.Url))

}

func (s *Content) GetPostsByUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userId := vars["id"]
	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		s.l.Println("converting id")
		return
	}

	stream, err := s.cc.GetPostsByUser(context.Background(), &prcontent.GetPostsRequest{UserId: id})
	if err != nil {
		s.l.Errorf("failed fetching posts %v\n", err)
		http.Error(w, "failed fetching posts", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("{"))
	for {
		posts, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.l.Errorf("failed to fetch posts: %v\n", err)
			http.Error(w, "Error couldn't fetch posts", http.StatusInternalServerError)
			return
		}
		saltdata.ToJSON(posts, w)
	}
	w.Write([]byte("}"))
}

func (s *Content) GetPostsByUserReaction(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userId := vars["id"]
	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		s.l.Println("converting id")
		return
	}

	stream, err := s.cc.GetPostsByUserReaction(context.Background(), &prcontent.GetPostsRequest{UserId: id})
	if err != nil {
		s.l.Errorf("failed fetching posts %v\n", err)
		http.Error(w, "failed fetching posts", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("{"))
	for {
		posts, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.l.Errorf("failed to fetch posts: %v\n", err)
			http.Error(w, "Error couldn't fetch posts", http.StatusInternalServerError)
			return
		}
		saltdata.ToJSON(posts, w)
	}
	w.Write([]byte("}"))
}

func (u *Users) GetFollowersDetailed(w http.ResponseWriter, r *http.Request) {

	jws, err := getUserJWS(r)
	if err != nil {
		u.l.Println("getting jws: %v\n", err)
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
		u.l.Errorf("parsing claims: %v", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		u.l.Errorf("unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	username, er := vars["username"]
	if !er {
		u.l.Println("[ERROR] parsing URL, no username in URL")
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	stream, err := u.uc.GetFollowersDetailed(context.Background(), &prusers.ProflieFollowRequest{Logeduser: claims.Username, Username: username})
	if err != nil {
		u.l.Println("[ERROR] fetching followers")
		http.Error(w, "Followers fetching error", http.StatusInternalServerError)
		return
	}
	var profiles []saltdata.ProfileFollowDetailedDTO
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
		dto := saltdata.ProfileFollowDetailedDTO{
			Username:  profile.Username,
			Following: profile.Following,
			Pending:   profile.Pending,
		}
		profiles = append(profiles, dto)
	}
	saltdata.ToJSON(profiles, w)
}

func (u *Users) GetFollowingDetailed(w http.ResponseWriter, r *http.Request) {

	jws, err := getUserJWS(r)
	if err != nil {
		u.l.Println("getting jws: %v\n", err)
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
		u.l.Errorf("parsing claims: %v", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		u.l.Errorf("unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	username, er := vars["username"]
	if !er {
		u.l.Println("[ERROR] parsing URL, no username in URL")
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	stream, err := u.uc.GetFollowingDetailed(context.Background(), &prusers.ProflieFollowRequest{Logeduser: claims.Username, Username: username})
	if err != nil {
		u.l.Println("[ERROR] fetching following")
		http.Error(w, "Followers fetching error", http.StatusInternalServerError)
		return
	}
	var profiles []saltdata.ProfileFollowDetailedDTO
	for {
		profile, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			u.l.Println("[ERROR] fetching followers")
			http.Error(w, "Error couldn't fetch following", http.StatusInternalServerError)
			return
		}
		dto := saltdata.ProfileFollowDetailedDTO{
			Username:  profile.Username,
			Following: profile.Following,
			Pending:   profile.Pending,
		}
		profiles = append(profiles, dto)
	}
	saltdata.ToJSON(profiles, w)
}

func (u *Users) CheckFollowing(w http.ResponseWriter, r *http.Request) {
	jws, err := getUserJWS(r)
	if err != nil {
		u.l.Println("getting jws: %v\n", err)
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
		u.l.Errorf("parsing claims: %v", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		u.l.Errorf("unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	username, er := vars["username"]
	if !er {
		u.l.Println("[ERROR] parsing URL, no username in URL")
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	resp, err := u.uc.CheckIfFollowing(context.Background(), &prusers.ProflieFollowRequest{Logeduser: claims.Username, Username: username})
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	saltdata.ToJSON(resp.Resposne, w)
}

func (u *Users) CheckFollowRequest(w http.ResponseWriter, r *http.Request) {
	jws, err := getUserJWS(r)
	if err != nil {
		u.l.Println("getting jws: %v\n", err)
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
		u.l.Errorf("parsing claims: %v", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		u.l.Errorf("unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	username, er := vars["username"]
	if !er {
		u.l.Println("[ERROR] parsing URL, no username in URL")
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	resp, err := u.uc.CheckForFollowingRequest(context.Background(), &prusers.ProflieFollowRequest{Logeduser: claims.Username, Username: username})
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	saltdata.ToJSON(resp.Resposne, w)
}
