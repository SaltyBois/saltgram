package handlers

import (
	"context"
	"io"
	"net/http"
	"os"
	"saltgram/data"
	saltdata "saltgram/data"
	"saltgram/protos/admin/pradmin"
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

	profile, err := u.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: user.Username, Username: user.Username})
	if err != nil {
		u.l.Println("[ERROR] fetching profile")
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	response := saltdata.UserDTO{
		Id:                strconv.FormatUint(user.Id, 10),
		Email:             user.Email,
		FullName:          user.FullName,
		Username:          user.Username,
		ProfilePictureURL: profile.ProfilePictureURL,
	}

	err = saltdata.ToJSON(response, w)
	if err != nil {
		u.l.Errorf("Serializing user: %v", err)
		http.Error(w, "Error serializing user", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
}

func (u *Users) GetByUsernameRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	user, err := u.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: username})
	if err != nil {
		u.l.Errorf("failed to fetch user: %v\n", err)
		http.Error(w, "User not found", http.StatusNotFound)
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
		u.l.Errorf("Serializing user: %v", err)
		http.Error(w, "Error serializing user", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
}

func (u *Users) GetProfile(w http.ResponseWriter, r *http.Request) {

	//jws, err := getUserJWS(r)
	//if err != nil {
	//	u.l.Println("[ERROR] JWS not found")
	//	http.Error(w, "JWS not found", http.StatusBadRequest)
	//	return
	//}
	//
	//token, err := jwt.ParseWithClaims(
	//	jws,
	//	&saltdata.AccessClaims{},
	//	func(t *jwt.Token) (interface{}, error) {
	//		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	//	},
	//)
	//
	//if err != nil {
	//	u.l.Printf("[ERROR] parsing claims: %v", err)
	//	http.Error(w, "Error parsing claims", http.StatusBadRequest)
	//	return
	//}
	//
	//claims, ok := token.Claims.(*saltdata.AccessClaims)
	//
	//if !ok {
	//	u.l.Println("[ERROR] unable to parse claims")
	//	http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
	//	return
	//}

	//user := claims.Username
	var isThisME bool

	userByJWS, err := getUserByJWS(r, u.uc)
	if err != nil {
		isThisME = false
	}

	vars := mux.Vars(r)
	profileUsername, er := vars["username"]
	if !er {
		u.l.Println("[ERROR] parsing URL, no username in URL")
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	if userByJWS != nil {
		isThisME = userByJWS.Username == profileUsername
		if !isThisME {
			blocked, _ := u.uc.CheckIfBlocked(context.Background(), &prusers.BlockProfileRequest{Logged: profileUsername, Profile: userByJWS.Username})
			if blocked.Response {
				http.Error(w, "Profile not found", http.StatusForbidden)
				return
			}
		}
	}

	profile, err := u.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: profileUsername, Username: profileUsername})
	if err != nil {
		u.l.Println("[ERROR] fetching profile")
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	response := saltdata.ProfileDTO{
		UserId:            strconv.FormatUint(profile.UserId, 10),
		FullName:          profile.FullName,
		Username:          profile.Username,
		Followers:         profile.Followers,
		Following:         profile.Following,
		Description:       profile.Description,
		IsPublic:          profile.IsPublic,
		IsFollowing:       profile.IsFollowing,
		PhoneNumber:       profile.PhoneNumber,
		Gender:            profile.Gender,
		DateOfBirth:       profile.DateOfBirth,
		WebSite:           profile.WebSite,
		ProfilePictureURL: profile.ProfilePictureURL,
		Taggable:          profile.Taggable,
		Messageable:       profile.Messageable,
		Verified:          profile.Verified,
		AccountType:       profile.AccountType,
		IsThisMe:          isThisME,
	}

	saltdata.ToJSON(response, w)
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

func (c *Content) GetHighlights(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["id"]
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		c.l.Errorf("failed to parse user id: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	resp, err := c.cc.GetHighlights(context.Background(), &prcontent.GetHighlightsRequest{UserId: userId})
	if err != nil {
		c.l.Errorf("failed to get highlights: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	dto := []*saltdata.HighlightDTO{}
	for _, h := range resp.Highlights {
		dto = append(dto, saltdata.PRToDTOHighlight(h))
	}
	saltdata.ToJSON(dto, w)
}

func (c *Content) GetStoriesByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["id"]
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		c.l.Errorf("failed to parse user id: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	resp, err := c.cc.GetStoriesIndividual(context.Background(), &prcontent.GetStoriesIndividualRequest{UserId: userId})
	if err != nil {
		c.l.Errorf("failed to get user stories: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	dto := []*saltdata.StoryDTO{}
	for _, s := range resp.Stories {
		dto = append(dto, saltdata.PRToDTOStory(s))
	}
	saltdata.ToJSON(dto, w)
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

	jws, err1 := getUserJWS(r)
	var t *jwt.Token
	var err2 error
	var claims1 *saltdata.AccessClaims
	var ok bool
	//if err != nil {
	//	u.l.Println("[ERROR] JWS not found")
	//	http.Error(w, "JWS not found", http.StatusBadRequest)
	//	return
	//}
	if jws != "" {
		token, error2 := jwt.ParseWithClaims(
			jws,
			&saltdata.AccessClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			},
		)
		t = token
		err2 = error2
	}

	//if err != nil {
	//	u.l.Printf("[ERROR] parsing claims: %v", err)
	//	http.Error(w, "Error parsing claims", http.StatusBadRequest)
	//	return
	//}
	if t != nil {
		claims, ok1 := t.Claims.(*saltdata.AccessClaims)
		claims1 = claims
		ok = ok1
	}

	//if !ok {
	//	u.l.Println("[ERROR] unable to parse claims")
	//	http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
	//	return
	//}

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
		if err1 == nil &&
			err2 == nil &&
			ok {
			if su.Username == claims1.Username {
				continue
			}
			blocked, _ := u.uc.CheckIfBlocked(context.Background(), &prusers.BlockProfileRequest{Logged: su.Username, Profile: claims1.Username})
			if blocked.Response {
				continue
			}
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

	type Message struct {
		Post        *prcontent.Post `json:"post"`
		User        *data.UserDTO   `json:"user"`
		TaggedUsers []*data.UserDTO `json:"taggedUsers"`
	}

	retVal := []*Message{}

	postsArray := []*prcontent.Post{}

	for {
		post, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.l.Errorf("failed to fetch posts: %v\n", err)
			http.Error(w, "Error couldn't fetch posts", http.StatusInternalServerError)
			return
		}
		postsArray = append(postsArray, post.Post)
	}

	for i := 0; i < len(postsArray); i++ {
		vr := postsArray[i]

		i, err := strconv.ParseUint(vr.UserId, 10, 64)
		if err != nil {
			s.l.Errorf("failed to convert id %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		user, err := s.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: i})

		if err != nil {
			s.l.Errorf("failed fetching user %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		profile, err := s.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: user.Username, Username: user.Username})
		if err != nil {
			s.l.Errorf("failed fetching profile %v\n", err)
			http.Error(w, "Profile getting error", http.StatusInternalServerError)
			return
		}

		taggedUsers := []*data.UserDTO{}
		for _, userTag := range vr.SharedMedia.Media[0].UserTags {
			userTagged, err := s.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: userTag.Id})

			if err != nil {
				s.l.Errorf("failed fetching user %v\n", err)
				http.Error(w, "User getting error", http.StatusInternalServerError)
				return
			}

			profileTagged, err := s.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: userTagged.Username, Username: userTagged.Username})
			if err != nil {
				s.l.Errorf("failed fetching profile %v\n", err)
				http.Error(w, "Profile getting error", http.StatusInternalServerError)
				return
			}

			taggedUsers = append(taggedUsers, &saltdata.UserDTO{
				Username:          profileTagged.Username,
				ProfilePictureURL: profileTagged.ProfilePictureURL,
			})
		}

		retVal = append(retVal, &Message{
			Post:        vr,
			TaggedUsers: taggedUsers,
			User: &data.UserDTO{
				Username:          profile.Username,
				ProfilePictureURL: profile.ProfilePictureURL,
			},
		})
	}
	saltdata.ToJSON(retVal, w)
}

/*func (s *Content) GetPostsByUserReaction(w http.ResponseWriter, r *http.Request) {

	jws, err := getUserJWS(r)
	if err != nil {
		s.l.Errorf("JWS not found: %v\n", err)
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
		s.l.Errorf("failure parsing claims: %v\n", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		s.l.Error("failed to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	user, err := s.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: claims.Username})
	if err != nil {
		s.l.Errorf("failed fetching user: %v\n", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	stream, err := s.cc.GetPostsByUserReaction(context.Background(), &prcontent.GetPostsRequest{UserId: user.Id})
	if err != nil {
		s.l.Errorf("failed fetching posts %v\n", err)
		http.Error(w, "failed fetching posts", http.StatusInternalServerError)
		return
	}
	var posts []*prcontent.GetPostsResponse
	for {
		post, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.l.Errorf("failed to fetch posts: %v\n", err)
			http.Error(w, "Error couldn't fetch posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}
	saltdata.ToJSON(posts, w)
}*/

func (s *Content) GetPostsByUserReaction(w http.ResponseWriter, r *http.Request) {

	jws, err := getUserJWS(r)
	if err != nil {
		s.l.Errorf("JWS not found: %v\n", err)
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
		s.l.Errorf("failure parsing claims: %v\n", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		s.l.Error("failed to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	user, err := s.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: claims.Username})
	if err != nil {
		s.l.Errorf("failed fetching user: %v\n", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	posts, err := s.cc.GetPostsByUserReaction(context.Background(), &prcontent.GetPostsByUserReactionRequest{Id: user.Id})
	if err != nil {
		s.l.Errorf("failed fetching posts %v\n", err)
		http.Error(w, "failed fetching posts", http.StatusInternalServerError)
		return
	}

	type Message struct {
		Post *prcontent.Post `json:"post"`
		User *data.UserDTO   `json:"user"`
	}

	postsArray := []*Message{}

	for i := 0; i < len(posts.Post); i++ {
		vr := posts.Post[i]

		i, err := strconv.ParseUint(vr.UserId, 10, 64)
		if err != nil {
			s.l.Errorf("failed to convert id %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		user, err := s.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: i})

		if err != nil {
			s.l.Errorf("failed fetching user %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		profile, err := s.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: user.Username, Username: user.Username})
		if err != nil {
			s.l.Errorf("failed fetching profile %v\n", err)
			http.Error(w, "Profile getting error", http.StatusInternalServerError)
			return
		}

		postsArray = append(postsArray, &Message{
			Post: vr,
			User: &data.UserDTO{
				Username:          user.Username,
				ProfilePictureURL: profile.ProfilePictureURL,
			},
		})
	}
	saltdata.ToJSON(postsArray, w)
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
			Username:       profile.Username,
			Following:      profile.Following,
			Pending:        profile.Pending,
			ProfliePicture: profile.ProfliePicture,
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

		user, err := u.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: profile.Username})
		if err != nil {
			u.l.Errorf("failed to fetch user: %v\n", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		dto := saltdata.ProfileFollowDetailedDTO{
			Username:       profile.Username,
			Following:      profile.Following,
			Pending:        profile.Pending,
			ProfliePicture: profile.ProfliePicture,
			Id:             strconv.FormatUint(user.Id, 10),
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

	saltdata.ToJSON(resp.Response, w)
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

	saltdata.ToJSON(resp.Response, w)
}
func (a *Admin) GetPendingVerifications(w http.ResponseWriter, r *http.Request) {

	verificationRequests, err := a.ac.GetPendingVerifications(context.Background(), &pradmin.GetVerificationRequest{})
	if err != nil {
		a.l.Errorf("failed fetching pending verifications %v\n", err)
		http.Error(w, "Pending verifications error", http.StatusInternalServerError)
		return
	}

	requests := []saltdata.VerificationRequestDTO{}

	for i := 0; i < len(verificationRequests.VerificationRequest); i++ {
		vr := verificationRequests.VerificationRequest[i]

		user, err := a.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: vr.UserId})

		if err != nil {
			a.l.Errorf("failed fetching user %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		profile, err := a.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: user.Username, Username: user.Username})
		if err != nil {
			a.l.Errorf("failed fetching profile %v\n", err)
			http.Error(w, "Profile getting error", http.StatusInternalServerError)
			return
		}

		requests = append(requests, saltdata.VerificationRequestDTO{

			Id:             strconv.FormatUint(vr.Id, 10),
			FullName:       vr.FullName,
			Category:       vr.Category,
			Url:            vr.Url,
			UserId:         vr.UserId,
			Username:       user.Username,
			ProfilePicture: profile.ProfilePictureURL,
		})
	}

	saltdata.ToJSON(&requests, w)
}

/*func (s *Content) GetCommentsByPost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userId := vars["id"]
	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		s.l.Println("converting id")
		return
	}

	stream, err := s.cc.GetComments(context.Background(), &prcontent.GetCommentsRequest{PostId: id})
	if err != nil {
		s.l.Errorf("failed fetching comments %v\n", err)
		http.Error(w, "failed fetching comments", http.StatusInternalServerError)
		return
	}
	var comments []*prcontent.GetCommentsResponse
	for {
		comment, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.l.Errorf("failed to fetch comments: %v\n", err)
			http.Error(w, "Error couldn't fetch comments", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}
	saltdata.ToJSON(comments, w)
}*/

func (s *Content) GetCommentsByPost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userId := vars["id"]
	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		s.l.Println("converting id")
		return
	}

	comments, err := s.cc.GetComments(context.Background(), &prcontent.GetCommentsRequest{PostId: id})
	if err != nil {
		s.l.Errorf("failed fetching comments %v\n", err)
		http.Error(w, "failed fetching comments", http.StatusInternalServerError)
		return
	}
	retVal := []saltdata.GetCommentDTO{}

	for i := 0; i < len(comments.Comment); i++ {
		vr := comments.Comment[i]

		user, err := s.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: vr.UserId})

		if err != nil {
			s.l.Errorf("failed fetching user %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		profile, err := s.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: user.Username, Username: user.Username})
		if err != nil {
			s.l.Errorf("failed fetching profile %v\n", err)
			http.Error(w, "Profile getting error", http.StatusInternalServerError)
			return
		}

		retVal = append(retVal, saltdata.GetCommentDTO{
			UserId:         strconv.FormatUint(vr.UserId, 10),
			Username:       user.Username,
			ProfilePicture: profile.ProfilePictureURL,
			Content:        vr.Content,
			PostId:         strconv.FormatUint(vr.PostId, 10),
		})
	}

	saltdata.ToJSON(&retVal, w)

}

func (s *Content) GetReactionsByPost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userId := vars["id"]
	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		s.l.Println("converting id")
		return
	}

	stream, err := s.cc.GetReactions(context.Background(), &prcontent.GetReactionsRequest{PostId: id})
	if err != nil {
		s.l.Errorf("failed fetching reactions %v\n", err)
		http.Error(w, "failed fetching reactions", http.StatusInternalServerError)
		return
	}
	var reactions []*prcontent.GetReactionsResponse
	for {
		reaction, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.l.Errorf("failed to fetch reactions: %v\n", err)
			http.Error(w, "Error couldn't fetch reactions", http.StatusInternalServerError)
			return
		}
		reactions = append(reactions, reaction)
	}
	saltdata.ToJSON(reactions, w)
}

func (a *Admin) GetPendingReports(w http.ResponseWriter, r *http.Request) {

	contentReports, err := a.ac.GetPendingInappropriateContentReport(context.Background(), &pradmin.GetInappropriateContentReportRequest{})
	if err != nil {
		a.l.Errorf("failed fetching pending reports %v\n", err)
		http.Error(w, "Pending reports error", http.StatusInternalServerError)
		return
	}

	reports := []saltdata.GetInappropriateContentReportDTO{}

	for i := 0; i < len(contentReports.InappropriateContentReport); i++ {
		vr := contentReports.InappropriateContentReport[i]

		user, err := a.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: vr.UserId})

		if err != nil {
			a.l.Errorf("failed fetching user %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		profile, err := a.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: user.Username, Username: user.Username})
		if err != nil {
			a.l.Errorf("failed fetching profile %v\n", err)
			http.Error(w, "Profile getting error", http.StatusInternalServerError)
			return
		}

		reports = append(reports, saltdata.GetInappropriateContentReportDTO{
			Id:             vr.Id,
			UserId:         vr.UserId,
			Username:       user.Username,
			ProfilePicture: profile.ProfilePictureURL,
			SharedMediaId:  vr.PostId,
			URL:            vr.Url,
		})
	}

	saltdata.ToJSON(&reports, w)
}

func (u *Users) GetMutedProfiles(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameByJWS(r)
	if err != nil {
		u.l.Println("failed to parse jws %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	stream, err := u.uc.GetMutedProfiles(context.Background(), &prusers.Profile{Username: username})
	if err != nil {
		u.l.Println("[ERROR] fetching muted profiles")
		http.Error(w, "Error fetching muted profiles", http.StatusInternalServerError)
		return
	}
	var profiles []*prusers.ProfileMBCF
	for {
		profile, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			u.l.Println("[ERROR] fetching muted profiles")
			http.Error(w, "Error couldn't fetch muted", http.StatusInternalServerError)
			return
		}

		profiles = append(profiles, profile)
	}
	saltdata.ToJSON(profiles, w)
}

func (u *Users) CheckIfMuted(w http.ResponseWriter, r *http.Request) {
	user, err := getUsernameByJWS(r)
	if err != nil {
		u.l.Println("failed to parse jws %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	username, er := vars["username"]
	if !er {
		u.l.Println("[ERROR] parsing URL, no username in URL")
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	resp, err := u.uc.CheckIfMuted(context.Background(), &prusers.MuteProfileRequest{Logged: user, Profile: username})
	if err != nil {
		http.Error(w, "Checking if muted", http.StatusNotFound)
		return
	}

	saltdata.ToJSON(resp.Response, w)

}

func (u *Users) GetBlockedProfiles(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameByJWS(r)
	if err != nil {
		u.l.Println("failed to parse jws %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	stream, err := u.uc.GetBlockedProfiles(context.Background(), &prusers.Profile{Username: username})
	if err != nil {
		u.l.Println("[ERROR] fetching blocked profiles")
		http.Error(w, "Error fetching blocked profiles", http.StatusInternalServerError)
		return
	}
	var profiles []*prusers.ProfileMBCF
	for {
		profile, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			u.l.Println("[ERROR] fetching followers")
			http.Error(w, "Error couldn't fetch blocked", http.StatusInternalServerError)
			return
		}
		profiles = append(profiles, profile)
	}
	saltdata.ToJSON(profiles, w)
}

func (u *Users) CheckIfBlocked(w http.ResponseWriter, r *http.Request) {
	user, err := getUsernameByJWS(r)
	if err != nil {
		u.l.Println("failed to parse jws %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	username, er := vars["username"]
	if !er {
		u.l.Println("[ERROR] parsing URL, no username in URL")
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	resp, err := u.uc.CheckIfBlocked(context.Background(), &prusers.BlockProfileRequest{Logged: user, Profile: username})
	if err != nil {
		http.Error(w, "Checking if blocked", http.StatusNotFound)
		return
	}

	saltdata.ToJSON(resp.Response, w)

}

func (u *Users) CheckIsBlocked(w http.ResponseWriter, r *http.Request) {
	user, err := getUsernameByJWS(r)
	if err != nil {
		u.l.Println("failed to parse jws %v", err)
		saltdata.ToJSON(false, w)
		return
	}

	vars := mux.Vars(r)
	username, er := vars["username"]
	if !er {
		u.l.Println("[ERROR] parsing URL, no username in URL")
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	resp, err := u.uc.CheckIfBlocked(context.Background(), &prusers.BlockProfileRequest{Logged: username, Profile: user})
	if err != nil {
		http.Error(w, "Checking if blocked", http.StatusNotFound)
		return
	}

	saltdata.ToJSON(resp.Response, w)

}

func (u *Users) GetCloseFriends(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameByJWS(r)
	if err != nil {
		u.l.Println("failed to parse jws %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	stream, err := u.uc.GetCloseFriends(context.Background(), &prusers.Profile{Username: username})
	if err != nil {
		u.l.Println("[ERROR] fetching close friends")
		http.Error(w, "Error fetching close friends", http.StatusInternalServerError)
		return
	}
	var profiles []*prusers.ProfileMBCF
	for {
		profile, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			u.l.Println("[ERROR] fetching profiles")
			http.Error(w, "Error couldn't fetch profile", http.StatusInternalServerError)
			return
		}

		profiles = append(profiles, profile)
	}
	saltdata.ToJSON(profiles, w)
}

func (u *Users) GetProfilesForCloseFriends(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameByJWS(r)
	if err != nil {
		u.l.Println("failed to parse jws %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	stream, err := u.uc.GetProfilesForCloseFriends(context.Background(), &prusers.Profile{Username: username})
	if err != nil {
		u.l.Println("[ERROR] fetching profiles for close friends")
		http.Error(w, "Error fetching profiles for close friends", http.StatusInternalServerError)
		return
	}
	var profiles []*prusers.ProfileMBCF
	for {
		profile, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			u.l.Println("[ERROR] fetching profile")
			http.Error(w, "Error couldn't fetch profile", http.StatusInternalServerError)
			return
		}

		profiles = append(profiles, profile)
	}
	saltdata.ToJSON(profiles, w)
}

func (s *Content) GetPostsByTag(w http.ResponseWriter, r *http.Request) {

	jws, err1 := getUserJWS(r)
	s.l.Errorf("error1 ", err1)

	var t *jwt.Token
	var err2 error
	var claims1 *saltdata.AccessClaims
	var ok bool
	//if err != nil {
	//	u.l.Println("[ERROR] JWS not found")
	//	http.Error(w, "JWS not found", http.StatusBadRequest)
	//	return
	//}
	if jws != "" {
		token, error2 := jwt.ParseWithClaims(
			jws,
			&saltdata.AccessClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			},
		)
		t = token
		err2 = error2
		s.l.Errorf("error2 ", err2)
	}

	//if err != nil {
	//	u.l.Printf("[ERROR] parsing claims: %v", err)
	//	http.Error(w, "Error parsing claims", http.StatusBadRequest)
	//	return
	//}
	if t != nil {
		claims, ok1 := t.Claims.(*saltdata.AccessClaims)
		claims1 = claims
		ok = ok1
		s.l.Errorf("ok ", ok)
	}

	vars := mux.Vars(r)
	value := vars["value"]

	posts, err := s.cc.SearchContent(context.Background(), &prcontent.SearchContentRequest{Value: value})
	if err != nil {
		s.l.Errorf("failed fetching posts %v\n", err)
		http.Error(w, "failed fetching posts", http.StatusInternalServerError)
		return
	}

	isUserLogged := false

	if err1 == nil && err2 == nil && ok {
		s.l.Errorf("logged in")
		isUserLogged = true
	} else {
		s.l.Errorf("not logged in")
	}

	type Message struct {
		Post *prcontent.Post `json:"post"`
		User *data.UserDTO   `json:"user"`
	}

	postsArray := []*Message{}

	for i := 0; i < len(posts.Post); i++ {
		vr := posts.Post[i]

		i, err := strconv.ParseUint(vr.UserId, 10, 64)
		if err != nil {
			s.l.Errorf("failed to convert id %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		user, err := s.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: i})

		if err != nil {
			s.l.Errorf("failed fetching user %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		profile, err := s.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: user.Username, Username: user.Username})
		if err != nil {
			s.l.Errorf("failed fetching profile %v\n", err)
			http.Error(w, "Profile getting error", http.StatusInternalServerError)
			return
		}

		isFollowing := false
		if isUserLogged {
			s.l.Errorf("My username", claims1.Username)
			following, err := s.uc.GerFollowing(context.Background(), &prusers.FollowerRequest{Username: claims1.Username})
			if err != nil {
				s.l.Println("[ERROR] fetching following", err)
				http.Error(w, "Following fetching error", http.StatusInternalServerError)
				return
			}
			for {
				prof, err := following.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					s.l.Println("[ERROR] fetching following", err)
					http.Error(w, "Error couldn't fetch following", http.StatusInternalServerError)
					return
				}
				s.l.Errorf("error             ", prof.Username)
				if prof.Username == profile.Username {
					isFollowing = true
				}
			}
		}
		isMyPost := false
		if err1 == nil && err2 == nil && ok {
			if profile.Username == claims1.Username {
				isMyPost = true
			}

		}

		if profile.IsPublic || isFollowing || isMyPost {
			postsArray = append(postsArray, &Message{
				Post: vr,
				User: &data.UserDTO{
					Username:          user.Username,
					ProfilePictureURL: profile.ProfilePictureURL,
				},
			})
		}
	}
	saltdata.ToJSON(postsArray, w)
}

func (c *Content) GetTagsByName(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	queryTag := vars["value"]

	queryResults, err := c.cc.GetTagsByName(context.Background(), &prcontent.GetTagsByNameRequest{Query: queryTag})

	if err != nil {
		c.l.Println("[ERROR] Searching tags failed")
		http.Error(w, "Error Searching tags: ", http.StatusInternalServerError)
		return
	}

	var finalResult []string

	const MAX_NUMBER_OF_RESULTS = 20

	for i := 0; i < len(queryResults.Name); i++ {
		if i == MAX_NUMBER_OF_RESULTS {
			break
		}
		finalResult = append(finalResult, queryResults.Name[i])
	}

	err = saltdata.ToJSON(&finalResult, w)

	if err != nil {
		c.l.Println("[ERROR] Searching tags failed")
		http.Error(w, "Error Searching tags: ", http.StatusInternalServerError)
		return
	}
}

func (c *Content) GetLocationNames(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	queryTag := vars["name"]

	queryResults, err := c.cc.GetLocationNames(context.Background(), &prcontent.GetLocationNamesRequest{Query: queryTag})

	if err != nil {
		c.l.Println("[ERROR] Searching location names failed")
		http.Error(w, "Error Searching location names: ", http.StatusInternalServerError)
		return
	}

	var finalResult []string

	const MAX_NUMBER_OF_RESULTS = 20

	for i := 0; i < len(queryResults.Name); i++ {
		if i == MAX_NUMBER_OF_RESULTS {
			break
		}
		finalResult = append(finalResult, queryResults.Name[i])
	}

	err = saltdata.ToJSON(&finalResult, w)

	if err != nil {
		c.l.Println("[ERROR] Searching location names failed")
		http.Error(w, "Error Searching location names: ", http.StatusInternalServerError)
		return
	}
}

func (s *Content) GetContentsByLocation(w http.ResponseWriter, r *http.Request) {

	jws, err1 := getUserJWS(r)
	s.l.Errorf("error1 ", err1)

	var t *jwt.Token
	var err2 error
	var claims1 *saltdata.AccessClaims
	var ok bool
	//if err != nil {
	//	u.l.Println("[ERROR] JWS not found")
	//	http.Error(w, "JWS not found", http.StatusBadRequest)
	//	return
	//}
	if jws != "" {
		token, error2 := jwt.ParseWithClaims(
			jws,
			&saltdata.AccessClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			},
		)
		t = token
		err2 = error2
		s.l.Errorf("error2 ", err2)
	}

	//if err != nil {
	//	u.l.Printf("[ERROR] parsing claims: %v", err)
	//	http.Error(w, "Error parsing claims", http.StatusBadRequest)
	//	return
	//}
	if t != nil {
		claims, ok1 := t.Claims.(*saltdata.AccessClaims)
		claims1 = claims
		ok = ok1
		s.l.Errorf("ok ", ok)
	}

	vars := mux.Vars(r)
	name := vars["name"]

	posts, err := s.cc.SearchContentLocation(context.Background(), &prcontent.SearchContentLocationRequest{Name: name})
	if err != nil {
		s.l.Errorf("failed fetching posts %v\n", err)
		http.Error(w, "failed fetching posts", http.StatusInternalServerError)
		return
	}

	isUserLogged := false

	if err1 == nil && err2 == nil && ok {
		s.l.Errorf("logged in")
		isUserLogged = true
	} else {
		s.l.Errorf("not logged in")
	}

	type Message struct {
		Post *prcontent.Post `json:"post"`
		User *data.UserDTO   `json:"user"`
	}

	postsArray := []*Message{}

	for i := 0; i < len(posts.Post); i++ {
		vr := posts.Post[i]

		i, err := strconv.ParseUint(vr.UserId, 10, 64)
		if err != nil {
			s.l.Errorf("failed to convert id %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		user, err := s.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: i})

		if err != nil {
			s.l.Errorf("failed fetching user %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		profile, err := s.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: user.Username, Username: user.Username})
		if err != nil {
			s.l.Errorf("failed fetching profile %v\n", err)
			http.Error(w, "Profile getting error", http.StatusInternalServerError)
			return
		}

		isFollowing := false
		if isUserLogged {
			following, err := s.uc.GerFollowing(context.Background(), &prusers.FollowerRequest{Username: claims1.Username})
			if err != nil {
				s.l.Println("[ERROR] fetching following", err)
				http.Error(w, "Following fetching error", http.StatusInternalServerError)
				return
			}
			for {
				prof, err := following.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					s.l.Println("[ERROR] fetching following", err)
					http.Error(w, "Error couldn't fetch following", http.StatusInternalServerError)
					return
				}
				if prof.Username == profile.Username {
					isFollowing = true
				}
			}
		}
		isMyPost := false
		if err1 == nil && err2 == nil && ok {
			if profile.Username == claims1.Username {
				isMyPost = true
			}
		}

		if profile.IsPublic || isFollowing || isMyPost {
			postsArray = append(postsArray, &Message{
				Post: vr,
				User: &data.UserDTO{
					Username:          user.Username,
					ProfilePictureURL: profile.ProfilePictureURL,
				},
			})
		}
	}
	saltdata.ToJSON(postsArray, w)
}

func (s *Users) GetTaggableProfiles(w http.ResponseWriter, r *http.Request) {

	myProfile, err := getProfileByJWS(r, s.uc)
	if err != nil {
		s.l.Println("getting jws: %v\n", err)
		http.Error(w, "JWS not found", http.StatusBadRequest)
		return
	}

	taggableProfiles := []*data.ProfileDTO{}

	following, err := s.uc.GerFollowing(context.Background(), &prusers.FollowerRequest{Username: myProfile.Username})
	if err != nil {
		s.l.Println("[ERROR] fetching following", err)
		http.Error(w, "Following fetching error", http.StatusInternalServerError)
		return
	}
	for {
		profile, err := following.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.l.Println("[ERROR] fetching following", err)
			http.Error(w, "Error couldn't fetch following", http.StatusInternalServerError)
			return
		}
		if profile.Taggable {
			s.l.Infof("username:   ", profile.Username)
			taggableProfiles = append(taggableProfiles, &data.ProfileDTO{
				Username:          profile.Username,
				ProfilePictureURL: profile.ProfilePicture,
				UserId:            profile.UserId,
			})
		}
	}

	saltdata.ToJSON(taggableProfiles, w)
}

func (s *Content) GetSavedPosts(w http.ResponseWriter, r *http.Request) {

	myProfile, err := getProfileByJWS(r, s.uc)
	if err != nil {
		s.l.Println("getting jws: %v\n", err)
		http.Error(w, "JWS not found", http.StatusBadRequest)
		return
	}

	posts, err := s.cc.GetSavedPosts(context.Background(), &prcontent.GetSavedPostsRequest{UserId: myProfile.UserId})
	if err != nil {
		s.l.Errorf("failed fetching posts %v\n", err)
		http.Error(w, "failed fetching posts", http.StatusInternalServerError)
		return
	}

	type Message struct {
		Post        *prcontent.Post `json:"post"`
		User        *data.UserDTO   `json:"user"`
		TaggedUsers []*data.UserDTO `json:"taggedUsers"`
	}

	postsArray := []*Message{}

	for i := 0; i < len(posts.Post); i++ {
		vr := posts.Post[i]

		i, err := strconv.ParseUint(vr.UserId, 10, 64)
		if err != nil {
			s.l.Errorf("failed to convert id %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		user, err := s.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: i})

		if err != nil {
			s.l.Errorf("failed fetching user %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		profile, err := s.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: user.Username, Username: user.Username})
		if err != nil {
			s.l.Errorf("failed fetching profile %v\n", err)
			http.Error(w, "Profile getting error", http.StatusInternalServerError)
			return
		}

		taggedUsers := []*data.UserDTO{}
		for _, userTag := range vr.SharedMedia.Media[0].UserTags {
			userTagged, err := s.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: userTag.Id})

			if err != nil {
				s.l.Errorf("failed fetching user %v\n", err)
				http.Error(w, "User getting error", http.StatusInternalServerError)
				return
			}

			profileTagged, err := s.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: userTagged.Username, Username: userTagged.Username})
			if err != nil {
				s.l.Errorf("failed fetching profile %v\n", err)
				http.Error(w, "Profile getting error", http.StatusInternalServerError)
				return
			}

			taggedUsers = append(taggedUsers, &saltdata.UserDTO{
				Username:          profileTagged.Username,
				ProfilePictureURL: profileTagged.ProfilePictureURL,
			})
		}

		postsArray = append(postsArray, &Message{
			Post: vr,
			User: &data.UserDTO{
				Username:          user.Username,
				ProfilePictureURL: profile.ProfilePictureURL,
			},
			TaggedUsers: taggedUsers,
		})
	}
	saltdata.ToJSON(postsArray, w)
}

func (s *Content) GetTaggedPosts(w http.ResponseWriter, r *http.Request) {

	myProfile, err := getProfileByJWS(r, s.uc)
	if err != nil {
		s.l.Println("getting jws: %v\n", err)
		http.Error(w, "JWS not found", http.StatusBadRequest)
		return
	}

	posts, err := s.cc.GetTaggedPosts(context.Background(), &prcontent.GetTaggedPostsRequest{UserId: myProfile.UserId})
	if err != nil {
		s.l.Errorf("failed fetching posts %v\n", err)
		http.Error(w, "failed fetching posts", http.StatusInternalServerError)
		return
	}

	type Message struct {
		Post        *prcontent.Post `json:"post"`
		User        *data.UserDTO   `json:"user"`
		TaggedUsers []*data.UserDTO `json:"taggedUsers"`
	}

	postsArray := []*Message{}

	for i := 0; i < len(posts.Post); i++ {
		vr := posts.Post[i]

		i, err := strconv.ParseUint(vr.UserId, 10, 64)
		if err != nil {
			s.l.Errorf("failed to convert id %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		user, err := s.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: i})

		if err != nil {
			s.l.Errorf("failed fetching user %v\n", err)
			http.Error(w, "User getting error", http.StatusInternalServerError)
			return
		}

		profile, err := s.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: user.Username, Username: user.Username})
		if err != nil {
			s.l.Errorf("failed fetching profile %v\n", err)
			http.Error(w, "Profile getting error", http.StatusInternalServerError)
			return
		}

		taggedUsers := []*data.UserDTO{}
		for _, userTag := range vr.SharedMedia.Media[0].UserTags {
			userTagged, err := s.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: userTag.Id})

			if err != nil {
				s.l.Errorf("failed fetching user %v\n", err)
				http.Error(w, "User getting error", http.StatusInternalServerError)
				return
			}

			profileTagged, err := s.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: userTagged.Username, Username: userTagged.Username})
			if err != nil {
				s.l.Errorf("failed fetching profile %v\n", err)
				http.Error(w, "Profile getting error", http.StatusInternalServerError)
				return
			}

			taggedUsers = append(taggedUsers, &saltdata.UserDTO{
				Username:          profileTagged.Username,
				ProfilePictureURL: profileTagged.ProfilePictureURL,
			})
		}

		postsArray = append(postsArray, &Message{
			Post: vr,
			User: &data.UserDTO{
				Username:          user.Username,
				ProfilePictureURL: profile.ProfilePictureURL,
			},
			TaggedUsers: taggedUsers,
		})
	}
	saltdata.ToJSON(postsArray, w)
}

func (u *Users) CheckActive(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, er := vars["username"]
	if !er {
		u.l.Println("[ERROR] parsing URL, no username in URL")
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	active, err := u.uc.CheckActive(context.Background(), &prusers.Profile{Username: username})
	if err != nil {
		u.l.Errorf("[ERROR] checking profile for username %v \n", err)
		return
	}

	saltdata.ToJSON(active.Response, w)
}

func (u *Users) GetFollowingMain(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameByJWS(r)
	if err != nil {
		u.l.Println("failed to parse jws %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	stream, err := u.uc.GetFollowingMain(context.Background(), &prusers.Profile{Username: username})
	if err != nil {
		u.l.Println("[ERROR] fetching following profiles")
		http.Error(w, "Error fetching following profiles", http.StatusInternalServerError)
		return
	}
	var profiles []saltdata.FollwingMainDTO
	for {
		profile, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			u.l.Println("[ERROR] fetching following")
			http.Error(w, "Error couldn't fetch following", http.StatusInternalServerError)
			return
		}

		user, err := u.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: profile.Username})
		if err != nil {
			u.l.Errorf("failed to fetch user: %v\n", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		dto := saltdata.FollwingMainDTO{
			Username:          profile.Username,
			ProfilePictureURL: profile.ProfilePictureURL,
			Id:                strconv.FormatUint(user.Id, 10),
		}
		profiles = append(profiles, dto)
	}
	saltdata.ToJSON(profiles, w)
}
