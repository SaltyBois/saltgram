package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	saltdata "saltgram/data"
	"saltgram/protos/admin/pradmin"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/content/prcontent"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
)

type ChangeRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

func (cr *ChangeRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(cr)
}

// func (e *Email) ChangePassword(w http.ResponseWriter, r *http.Request) {
// 	cr := ChangeRequest{}

// 	err := saltdata.FromJSON(&cr, r.Body)
// 	if err != nil {
// 		e.l.Printf("[ERROR] deserializing ChangeRequest: %v\n", err)
// 		http.Error(w, "Failed to parse request", http.StatusBadRequest)
// 		return
// 	}

// 	err = cr.Validate()
// 	if err != nil {
// 		e.l.Printf("[ERROR] ChangeRequest not valid: %v\n", err)
// 		http.Error(w, "Bad change request", http.StatusBadRequest)
// 		return
// 	}

// 	cookie, err := r.Cookie("email")
// 	if err != nil {
// 		e.l.Printf("[ERROR] getting cookie: %v", err)
// 		http.Error(w, "No change request cookie", http.StatusBadRequest)
// 		return
// 	}
// 	_, err = e.uc.ChangePassword(context.Background(), &prusers.ChangeRequest{
// 		Email:            cookie.Value,
// 		OldPlainPassword: cr.OldPassword,
// 		NewPlainPassword: cr.NewPassword,
// 	})
// 	if err != nil {
// 		e.l.Printf("[ERROR] POST change password request: %v\n", err)
// 		http.Error(w, "Error in POST change password request", http.StatusBadRequest)
// 		return
// 	}
// 	w.Write([]byte("200 - OK"))
// }

func (e *Email) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e.l.Errorf("failed to get email: %v\n", err)
		http.Error(w, "No email", http.StatusBadRequest)
		return
	}
	email := string(body)
	go func() {
		_, err = e.ec.RequestReset(context.Background(), &premail.ResetRequest{Email: email})
		if err != nil {
			e.l.Errorf("failed sending email request: %v\n", err)
		}
	}()
	// NOTE(Jovan): Always send 200 OK as per OWASP
	w.Write([]byte("200 - OK"))
}

func (u *Users) ChangePassword(w http.ResponseWriter, r *http.Request) {
	cr := ChangeRequest{}

	err := saltdata.FromJSON(&cr, r.Body)
	if err != nil {
		u.l.Errorf("failed to deserialize ChangeRequest: %v\n", err)
		http.Error(w, "Failed to parse request", http.StatusBadRequest)
		return
	}

	err = cr.Validate()
	if err != nil {
		u.l.Errorf("ChangeRequest not valid: %v\n", err)
		http.Error(w, "Bad change request", http.StatusBadRequest)
		return
	}
	jws, err := getUserJWS(r)
	if err != nil {
		u.l.Errorf("failed to get jws: %v\n", err)
		http.Error(w, "Missing jws", http.StatusBadRequest)
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
		u.l.Errorf("failed parsing claims: %v\n", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		u.l.Error("unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	_, err = u.uc.ChangePassword(context.Background(), &prusers.ChangeRequest{
		Username:         claims.Username,
		OldPlainPassword: cr.OldPassword,
		NewPlainPassword: cr.NewPassword,
	})

	if err != nil {
		u.l.Errorf("failed to change password: %v\n", err)
		http.Error(w, "Error in POST change password request", http.StatusBadRequest)
		return
	}
	w.Write([]byte("200 - OK"))
}

func (u *Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		u.l.Errorf("failure parsing body: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("emailforreset")
	if err != nil {
		u.l.Errorf("failed getting cookie: %v\n", err)
		http.Error(w, "No reset request cookie", http.StatusBadRequest)
		return
	}

	newPassword := string(body)
	_, err = u.uc.ResetPassword(context.Background(), &prusers.UserResetRequest{Email: cookie.Value, Password: newPassword})
	if err != nil {
		u.l.Errorf("failed resetting password: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Write([]byte("200 - OK"))
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	dto := saltdata.UserDTO{}
	err := saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		u.l.Errorf("failed deserializing user data: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = dto.Validate()
	if err != nil {
		u.l.Errorf("failed validating user data: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = u.uc.Register(context.Background(), &prusers.RegisterRequest{
		Username:    dto.Username,
		FullName:    dto.FullName,
		Email:       dto.Email,
		Password:    dto.Password,
		Description: dto.Description,
		ReCaptcha: &prusers.UserReCaptcha{
			Token:  dto.ReCaptcha.Token,
			Action: dto.ReCaptcha.Action,
		},
		PhoneNumber:    dto.PhoneNumber,
		Gender:         dto.Gender,
		DateOfBirth:    dto.DateOfBirth.Unix(),
		WebSite:        dto.WebSite,
		PrivateProfile: dto.PrivateProfile,
	})

	if err != nil {
		u.l.Errorf("failed registering user: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	u.l.Infof("User registered: %v\n", dto.Email)
	w.Write([]byte("Activation email sent"))
}

func (a *Auth) Get2FAQR(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.l.Errorf("failure reading request body: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	username := string(body)
	res, err := a.ac.Get2FAQR(context.Background(), &prauth.TwoFARequest{Username: username})
	if err != nil {
		a.l.Errorf("failed to get 2FAQR, returned: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(res.Png)
	if err != nil {
		a.l.Errorf("failed to write png data: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func (a *Auth) GetJWT(w http.ResponseWriter, r *http.Request) {
	user := saltdata.Login{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		a.l.Errorf("failure deserializing user: %v\n", err)
		http.Error(w, "Failed to deserialize user", http.StatusBadRequest)
		return
	}
	res, err := a.ac.GetJWT(context.Background(), &prauth.JWTRequest{Username: user.Username, Password: user.Password})
	if err != nil {
		a.l.Errorf("failure getting jwt: %v\n", err)
		http.Error(w, "Faile to get jwt", http.StatusBadRequest)
		return
	}
	cookie := http.Cookie{
		Name:     "refresh",
		Value:    res.Refresh,
		Expires:  time.Now().UTC().AddDate(0, 6, 0),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}

	http.SetCookie(w, &cookie)

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(res.Jws))
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	login := saltdata.Login{}
	err := saltdata.FromJSON(&login, r.Body)
	if err != nil {
		a.l.Errorf("failure deserializing body: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = login.Validate()
	if err != nil {
		a.l.Errorf("failure validating: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	res, err := a.ac.Login(context.Background(), &prauth.LoginRequest{
		Username: login.Username,
		Password: login.Password,
		ReCaptcha: &prauth.ReCaptcha{
			Action: login.ReCaptcha.Action,
			Token:  login.ReCaptcha.Token,
		},
	})

	if err != nil {
		a.l.Errorf("failure calling login: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	a.l.Infof("%v logged in\n", login.Username)
	saltdata.ToJSON(res, w)
}

// What?
// func (c *Content) AddSharedMedia(w http.ResponseWriter, r *http.Request) {

// 	jws, err := getUserJWS(r)
// 	if err != nil {
// 		c.l.Errorf("JWS not found: %v\n", err)
// 		http.Error(w, "JWS not found", http.StatusBadRequest)
// 		return
// 	}

// 	token, err := jwt.ParseWithClaims(
// 		jws,
// 		&saltdata.AccessClaims{},
// 		func(t *jwt.Token) (interface{}, error) {
// 			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
// 		},
// 	)

// 	if err != nil {
// 		c.l.Errorf("failure parsing claims: %v\n", err)
// 		http.Error(w, "Error parsing claims", http.StatusBadRequest)
// 		return
// 	}

// 	claims, ok := token.Claims.(*saltdata.AccessClaims)

// 	if !ok {
// 		c.l.Error("failed to parse claims")
// 		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
// 		return
// 	}

// 	user, err := c.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: claims.Username})
// 	if err != nil {
// 		c.l.Errorf("failed fetching user: %v\n", err)
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	dto := saltdata.SharedMediaDTO{}
// 	err = saltdata.FromJSON(&dto, r.Body)
// 	if err != nil {
// 		c.l.Errorf("failure adding shared media: %v\n", err)
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 		return
// 	}

// 	media := []*prcontent.Media{}
// 	for _, m := range dto.Media {
// 		tags := []*prcontent.Tag{}
// 		for _, t := range m.Tags {
// 			tags = append(tags, &prcontent.Tag{
// 				Value: t.Value,
// 				Id:    t.ID,
// 			})
// 		}
// 		media = append(media, &prcontent.Media{
// 			UserId:      user.Id,
// 			Filename:    m.Filename,
// 			Tags:        tags,
// 			Description: m.Description,
// 			Location: &prcontent.Location{
// 				Country: m.Location.Country,
// 				State:   m.Location.State,
// 				ZipCode: m.Location.ZipCode,
// 				Street:  m.Location.Street,
// 			},
// 			AddedOn: m.AddedOn,
// 		})
// 	}

// 	_, err = c.cc.AddSharedMedia(context.Background(), &prcontent.AddSharedMediaRequest{Media: media})
// 	if err != nil {
// 		c.l.Errorf("failed to add shared media: %v\n", err)
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 		return
// 	}
// }

func (c *Content) AddHighlight(w http.ResponseWriter, r *http.Request) {
	user, err := getUserByJWS(r, c.uc)
	if err != nil {
		c.l.Errorf("failed fetching user: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	dto := saltdata.HighlightRequest{}
	err = saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		c.l.Errorf("failed to parse request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	stories := []uint64{}
	for _, s := range dto.Stories {
		id, err := strconv.ParseUint(s.Id, 10, 64)
		if err != nil {
			c.l.Errorf("failed to parse uint: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		stories = append(stories, id)
	}
	_, err = c.cc.AddHighlight(context.Background(), &prcontent.AddHighlightRequest{Name: dto.Name, UserId: user.Id, Stories: stories})
	if err != nil {
		c.l.Errorf("failed to add highlight")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Added"))
}

func (c *Content) AddProfilePicture(w http.ResponseWriter, r *http.Request) {

	profile, err := getProfileByJWS(r, c.uc)
	if err != nil {
		c.l.Errorf("failed fetching user: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	c.l.Infof("fetched user: %v, id: %v", profile.Username, profile.UserId)

	// TODO(Jovan): Authenticate

	r.ParseMultipartForm(10 << 20) // up to 10MB
	file, _, err := r.FormFile("profileImg")
	if err != nil {
		c.l.Errorf("failed to parse request: %v", err)
		http.Error(w, "Invalid image data", http.StatusBadRequest)
	}
	defer file.Close()

	imageBytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.l.Errorf("failed to read image bytes: %v", err)
		http.Error(w, "Invalid image data", http.StatusBadRequest)
		return
	}

	stream, err := c.cc.AddProfilePicture(context.Background())
	if err != nil {
		c.l.Errorf("failed to add profile picture: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = stream.Send(&prcontent.AddProfilePictureRequest{
		Data: &prcontent.AddProfilePictureRequest_Info{Info: &prcontent.AddProfilePictureInfo{
			UserId:          profile.UserId,
			ProfileFolderId: profile.ProfileFolderId,
		}},
	})
	if err != nil {
		c.l.Errorf("failed to send profile picture metadata: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = stream.Send(&prcontent.AddProfilePictureRequest{
		Data: &prcontent.AddProfilePictureRequest_Image{Image: imageBytes},
	})
	if err != nil {
		c.l.Errorf("failed to send profile picture bytes: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		c.l.Errorf("failed to recieve picture url: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	c.l.Infof("received response: %v", resp.Url)
	_, err = c.uc.UpdateProfilePicture(context.Background(), &prusers.UpdateProfilePictureRequest{Url: resp.Url, Username: profile.Username})
	if err != nil {
		c.l.Errorf("failed to update profiel picture url: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(resp.Url))
}

func (c *Content) AddStory(w http.ResponseWriter, r *http.Request) {
	profile, err := getProfileByJWS(r, c.uc)
	if err != nil {
		c.l.Errorf("failed fetching user: %v\n", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	r.ParseMultipartForm(10 * 10 << 20) // Max 10 * 10MB
	formdata := r.MultipartForm
	files := formdata.File["stories"]
	if len(files) == 0 {
		c.l.Error("empty files form")
		http.Error(w, "Empty files form", http.StatusBadRequest)
		return
	}

	// TODO(Jovan): Tag value as ID or predefined tags
	tags := []*prcontent.Tag{}
	// TODO(Jovan): Pass location as object
	// location := r.PostForm["location"]
	description := ""
	if len(r.PostForm["description"]) > 0 {
		description = r.PostForm["description"][0]
	}

	// NOTE(Jovan): default
	location := &prcontent.Location{
		Country: "RS",
		State:   "Serbia",
		ZipCode: "21000",
		Street:  "Balzakova 69",
	}

	if len(r.PostForm["location"]) > 0 {
		err = json.Unmarshal([]byte(r.PostForm["location"][0]), &location)
		if err != nil {
			c.l.Errorf("failed to unmarshal location: %v", err)
			http.Error(w, "Failed to parse location", http.StatusBadRequest)
			return
		}
	}

	resp, err := c.cc.CreateStory(context.Background(), &prcontent.CreateStoryRequest{UserId: profile.UserId})
	if err != nil {
		c.l.Errorf("failed to create shared media: %v", err)
		http.Error(w, "Failed to create shared media", http.StatusInternalServerError)
		return
	}

	c.l.Infof("received %v files", len(files))

	for _, f := range files {
		media := &prcontent.Media{
			UserId:      profile.UserId,
			Filename:    f.Filename,
			Tags:        tags,
			Description: description,
			Location:    location,
			AddedOn:     time.Now().String(),
		}
		stream, err := c.cc.AddStory(context.Background())
		if err != nil {
			c.l.Errorf("failed to add story: %v", err)
			http.Error(w, "Failed to add story", http.StatusInternalServerError)
			return
		}

		closeFriends := false
		if len(r.PostForm["closeFriends"]) > 0 {
			c.l.Info("friends: %v", r.PostForm["closeFriends"])
			closeFriends = r.PostForm["closeFriends"][0] == "true"
		}

		err = stream.Send(&prcontent.AddStoryRequest{Data: &prcontent.AddStoryRequest_Info{
			Info: &prcontent.AddStoryRequestInfo{
				UserId:          profile.UserId,
				StoriesFolderId: profile.StoriesFolderId,
				CloseFriends:    closeFriends,
				Media:           media,
				StoryId:         resp.StoryId,
			},
		}})

		if err != nil {
			c.l.Errorf("failed to send story meta: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		file, err := f.Open()
		if err != nil {
			c.l.Errorf("failed to open file: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		imgBytes, err := ioutil.ReadAll(file)
		if err != nil {
			c.l.Errorf("failed to read file bytes: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		err = stream.Send(&prcontent.AddStoryRequest{Data: &prcontent.AddStoryRequest_Image{
			Image: imgBytes,
		}})

		if err != nil {
			c.l.Errorf("failed to send image data: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		_, err = stream.CloseAndRecv()
		if err != nil {
			c.l.Errorf("failed to close and recieve: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("Added"))
}

func (c *Content) AddPost(w http.ResponseWriter, r *http.Request) {

	profile, err := getProfileByJWS(r, c.uc)
	if err != nil {
		c.l.Errorf("failed fetching user: %v\n", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	r.ParseMultipartForm(10 * 10 << 20) // Max 10 * 10MB
	formdata := r.MultipartForm
	files := formdata.File["posts"]
	if len(files) == 0 {
		c.l.Error("empty files form")
		http.Error(w, "Empty files form", http.StatusBadRequest)
		return
	}

	// TODO(Jovan): Tag value as ID or predefined tags
	tags := []*prcontent.Tag{}
	// TODO(Jovan): Pass location as object
	// location := r.PostForm["location"]
	description := ""
	if len(r.PostForm["description"]) > 0 {
		description = r.PostForm["description"][0]
	}

	// NOTE(Jovan): default
	location := &prcontent.Location{
		Country: "RS",
		State:   "Serbia",
		ZipCode: "21000",
		Street:  "Balzakova 69",
	}

	if len(r.PostForm["location"]) > 0 {
		err = json.Unmarshal([]byte(r.PostForm["location"][0]), &location)
		if err != nil {
			c.l.Errorf("failed to unmarshal location: %v", err)
			http.Error(w, "Failed to parse location", http.StatusBadRequest)
			return
		}
	}

	resp, err := c.cc.CreatePost(context.Background(), &prcontent.CreatePostRequest{UserId: profile.UserId})
	if err != nil {
		c.l.Errorf("failed to create shared media: %v", err)
		http.Error(w, "Failed to create shared media", http.StatusInternalServerError)
		return
	}

	for _, f := range files {
		media := &prcontent.Media{
			UserId:      profile.UserId,
			Filename:    f.Filename,
			Tags:        tags,
			Description: description,
			Location:    location,
			AddedOn:     time.Now().String(),
		}
		stream, err := c.cc.AddPost(context.Background())
		if err != nil {
			c.l.Errorf("failed to add post: %v", err)
			http.Error(w, "Failed to add post", http.StatusInternalServerError)
			return
		}

		err = stream.Send(&prcontent.AddPostRequest{Data: &prcontent.AddPostRequest_Info{
			Info: &prcontent.AddPostRequestInfo{
				Media:         media,
				UserId:        profile.UserId,
				PostsFolderId: profile.PostsFolderId,
				PostId:        resp.PostId,
			},
		}})
		if err != nil {
			c.l.Errorf("failed to send post meta: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		file, err := f.Open()
		if err != nil {
			c.l.Errorf("failed to open file: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		imgBytes, err := ioutil.ReadAll(file)
		if err != nil {
			c.l.Errorf("failed to read file bytes: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		err = stream.Send(&prcontent.AddPostRequest{Data: &prcontent.AddPostRequest_Image{
			Image: imgBytes,
		}})

		if err != nil {
			c.l.Errorf("failed to send image data: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		_, err = stream.CloseAndRecv()
		if err != nil {
			c.l.Errorf("failed to close and recieve: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("Added"))
}

func (c *Content) AddComment(w http.ResponseWriter, r *http.Request) {

	jws, err := getUserJWS(r)
	if err != nil {
		c.l.Errorf("JWS not found: %v\n", err)
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
		c.l.Errorf("failure parsing claims: %v\n", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		c.l.Error("failed to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	user, err := c.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: claims.Username})
	if err != nil {
		c.l.Errorf("failed fetching user: %v\n", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	dto := saltdata.CommentDTO{}
	err = saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		c.l.Errorf("failure adding comment: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	i, err := strconv.ParseUint(dto.PostId, 0, 64)
	if err != nil {
		c.l.Errorf("failed to convert id post: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = c.cc.AddComment(context.Background(), &prcontent.AddCommentRequest{Content: dto.Content, UserId: user.Id, PostId: i})

	if err != nil {
		c.l.Errorf("failed to add post: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func (c *Content) AddReaction(w http.ResponseWriter, r *http.Request) {

	jws, err := getUserJWS(r)
	if err != nil {
		c.l.Errorf("JWS not found: %v\n", err)
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
		c.l.Errorf("failure parsing claims: %v\n", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		c.l.Error("failed to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	user, err := c.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: claims.Username})
	if err != nil {
		c.l.Errorf("failed fetching user: %v\n", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	dto := saltdata.ReactionDTO{}
	err = saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		c.l.Errorf("failure adding reaction: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	i, err := strconv.ParseUint(dto.PostId, 0, 64)
	if err != nil {
		c.l.Errorf("failed to convert id post: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = c.cc.AddReaction(context.Background(), &prcontent.AddReactionRequest{ReactionType: dto.ReactionType, UserId: user.Id, PostId: i})

	if err != nil {
		c.l.Errorf("failed to add reaction: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func (u *Users) Follow(w http.ResponseWriter, r *http.Request) {
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

	profileRequest := claims.Username

	dto := saltdata.FollowDTO{}
	err = saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		u.l.Printf("[ERROR] deserializing user data: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = dto.Validate()
	if err != nil {
		u.l.Printf("[ERROR] validating user data: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = u.uc.Follow(context.Background(), &prusers.FollowRequest{Username: profileRequest, ToFollow: dto.ProfileToFollow})
	if err != nil {
		u.l.Printf("[ERROR] following profile: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Following %s"))

}

func (a *Admin) AddVerificationRequest(w http.ResponseWriter, r *http.Request) {
	a.l.Info("Requesting verification")
	user, err := getUserByJWS(r, a.uc)
	if err != nil {
		a.l.Errorf("failed fetching user: %v\n", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	r.ParseMultipartForm(10 << 20) // up to 10MB
	formdata := r.MultipartForm
	files := formdata.File["document"]
	if len(files) == 0 {
		a.l.Errorf("empty file request")
		http.Error(w, "Invalid document data", http.StatusBadRequest)
	}

	file, err := files[0].Open()
	if err != nil {
		a.l.Errorf("failed to open document file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	imageBytes, err := ioutil.ReadAll(file)
	if err != nil {
		a.l.Errorf("failed to read document bytes: %v", err)
		http.Error(w, "Invalid document data", http.StatusBadRequest)
		return
	}
	category := r.FormValue("category")
	fullName := r.FormValue("fullName")

	stream, err := a.cc.AddVerificationImage(context.Background())
	if err != nil {
		a.l.Errorf("failed to open content stream: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = stream.Send(&prcontent.AddVerificationImageRequest{Data: &prcontent.AddVerificationImageRequest_Info{
		Info: &prcontent.AddVerificationImageRequestInfo{
			UserId:   user.Id,
			Filename: files[0].Filename,
		},
	}})
	if err != nil {
		a.l.Errorf("failed to send document meta: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = stream.Send(&prcontent.AddVerificationImageRequest{Data: &prcontent.AddVerificationImageRequest_Image{
		Image: imageBytes,
	}})
	if err != nil {
		a.l.Errorf("failed to send document data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		a.l.Errorf("failed to receive document url: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = a.ac.AddVerificationReq(context.Background(), &pradmin.AddVerificationRequest{
		FullName: fullName,
		UserId:   user.Id,
		Category: category,
		Url:      resp.Url,
	})

	if err != nil {
		a.l.Errorf("failed to add verification request: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Write([]byte(resp.Url))
}

func (a *Admin) SendInappropriateContentReport(w http.ResponseWriter, r *http.Request) {

	jws, err := getUserJWS(r)
	if err != nil {
		a.l.Errorf("JWS not found: %v\n", err)
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
		a.l.Errorf("failure parsing claims: %v\n", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		a.l.Error("failed to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	user, err := a.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: claims.Username})
	if err != nil {
		a.l.Errorf("failed fetching user: %v\n", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	dto := saltdata.InappropriateContentReportDTO{}
	err = saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		a.l.Errorf("failure adding inappropriate content: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	i, err := strconv.ParseUint(dto.SharedMediaId, 10, 64)
	if err != nil {
		a.l.Println("converting id")
		return
	}

	_, err = a.ac.SendInappropriateContentReport(context.Background(), &pradmin.InappropriateContentReportRequest{UserId: user.Id, PostId: i})

	if err != nil {
		a.l.Errorf("failed to add inappropriate: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}
