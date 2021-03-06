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
	"saltgram/protos/notifications/prnotifications"
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

func (u *Users) AcceptCampaign(w http.ResponseWriter, r *http.Request) {
	user, err := getUserByJWS(r, u.uc)
	if err != nil {
		u.l.Errorf("failed to get user by jws: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		u.l.Errorf("failure parsing body: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id, _ := strconv.ParseUint(string(body), 10, 64)

	_, err = u.uc.AcceptInfluencer(context.Background(), &prusers.AcceptInfluencerRequest{
		InfluencerId: user.Id,
		CampaignId:   id,
	})

	if err != nil {
		u.l.Errorf("failed to accept influencer request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	u.l.Infof("Adding influencer: %v", user.Id)

	_, err = u.cc.AddInfluencerToCampaign(context.Background(), &prcontent.AddInfluencerToCampaignRequest{
		InfluencerId: user.Id,
		CampaignId:   id,
	})

	if err != nil {
		u.l.Errorf("failed to add influencer to campaign: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Accepted"))
}

func (u *Users) SendInfluencerRequest(w http.ResponseWriter, r *http.Request) {
	_, err := getUserByJWS(r, u.uc)
	if err != nil {
		u.l.Errorf("failed to get user by jws: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	dto := struct {
		InfluencerID string `json:"influencerId"`
		CampaignID   string `json:"campaignId"`
		Website      string `json:"website"`
	}{}
	err = saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		u.l.Errorf("failed to read influencer request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	influencerId, _ := strconv.ParseUint(dto.InfluencerID, 10, 64)
	campaignId, _ := strconv.ParseUint(dto.CampaignID, 10, 64)
	_, err = u.uc.InfluencerRequest(context.Background(), &prusers.InfluencerRequestRequest{
		InfluencerId: influencerId,
		CampaignId:   campaignId,
		Website:      dto.Website,
	})
	if err != nil {
		u.l.Errorf("failed to add influencer request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.Write([]byte("Requested"))
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
		Agent:          dto.Agent,
	})

	if err != nil {
		u.l.Errorf("failed registering user: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if dto.Agent {
		_, err = u.ac.AddAgentRegistration(context.Background(), &pradmin.AddAgentRegistrationRequest{AgentEmail: dto.Email})
		if err != nil {
			u.l.Errorf("failed to add agent request: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
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

	tags := []*prcontent.Tag{}
	for _, t := range r.PostForm["tags"] {
		tags = append(tags, &prcontent.Tag{
			Value: t,
		})
	}
	userTags := []*prcontent.UserTag{}
	for _, t := range r.PostForm["userTags"] {
		i, err := strconv.ParseUint(t, 10, 64)
		if err != nil {
			c.l.Errorf("failed to parse user tag id: %v", err)
			http.Error(w, "Failed to parse user tag id", http.StatusBadRequest)
			return
		}
		userTags = append(userTags, &prcontent.UserTag{
			Id: i,
		})
	}
	description := ""
	if len(r.PostForm["description"]) > 0 {
		description = r.PostForm["description"][0]
	}

	location := &prcontent.Location{}

	if len(r.PostForm["location"]) > 0 {
		err = json.Unmarshal([]byte(r.PostForm["location"][0]), &location)
		if err != nil {
			c.l.Errorf("failed to unmarshal location: %v", err)
			http.Error(w, "Failed to parse location", http.StatusBadRequest)
			return
		}
	}

	isCampaign := false
	if len(r.PostForm["campaign"]) > 0 {
		err = json.Unmarshal([]byte(r.PostForm["campaign"][0]), &isCampaign)
		if err != nil {
			c.l.Errorf("failed to unmarshal campaign bool: %v", err)
			http.Error(w, "Failed to parse campaign bool", http.StatusBadRequest)
			return
		}
	}
	ageGroup := "Pre 20s"
	if len(r.PostForm["ageGroup"]) > 0 {
		ageGroup = r.PostForm["ageGroup"][0]
	}
	campaignWebsite := r.PostForm.Get("website")

	campaignOneTime := false
	if len(r.PostForm["oneTime"]) > 0 {
		err = json.Unmarshal([]byte(r.PostForm["oneTime"][0]), &campaignOneTime)
		if err != nil {
			c.l.Errorf("failed to unmarshal oneTime bool: %v", err)
			http.Error(w, "Failed to parse oneTime bool", http.StatusBadRequest)
			return
		}
	}
	campaignStart := r.PostForm.Get("campaignStart")
	campaignEnd := r.PostForm.Get("campaignEnd")

	resp, err := c.cc.CreateStory(context.Background(), &prcontent.CreateStoryRequest{
		UserId:          profile.UserId,
		Campaign:        isCampaign,
		AgeGroup:        ageGroup,
		CampaignOneTime: campaignOneTime,
		CampaignStart:   campaignStart,
		CampaignEnd:     campaignEnd,
		CampaignWebsite: campaignWebsite,
	})
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
			UserTags:    userTags,
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
	for _, t := range r.PostForm["tags"] {
		tags = append(tags, &prcontent.Tag{
			Value: t,
		})
	}

	userTags := []*prcontent.UserTag{}
	for _, t := range r.PostForm["userTags"] {
		i, err := strconv.ParseUint(t, 10, 64)
		if err != nil {
			c.l.Errorf("failed to parse user tag id: %v", err)
			http.Error(w, "Failed to parse user tag id", http.StatusBadRequest)
			return
		}
		userTags = append(userTags, &prcontent.UserTag{
			Id: i,
		})
	}
	// TODO(Jovan): Pass location as object
	// location := r.PostForm["location"]
	description := ""
	if len(r.PostForm["description"]) > 0 {
		description = r.PostForm["description"][0]
	}

	// NOTE(Jovan): default
	location := &prcontent.Location{
		/*Country: "RS",
		State:   "Serbia",
		ZipCode: "21000",
		Street:  "Balzakova 69",*/
	}

	if len(r.PostForm["location"]) > 0 {
		err = json.Unmarshal([]byte(r.PostForm["location"][0]), &location)
		if err != nil {
			c.l.Errorf("failed to unmarshal location: %v", err)
			http.Error(w, "Failed to parse location", http.StatusBadRequest)
			return
		}
	}

	isCampaign := false
	if len(r.PostForm["campaign"]) > 0 {
		err = json.Unmarshal([]byte(r.PostForm["campaign"][0]), &isCampaign)
		if err != nil {
			c.l.Errorf("failed to unmarshal campaign bool: %v", err)
			http.Error(w, "Failed to parse campaign bool", http.StatusBadRequest)
			return
		}
	}
	ageGroup := "Pre 20s"
	if len(r.PostForm["ageGroup"]) > 0 {
		ageGroup = r.PostForm["ageGroup"][0]
	}
	campaignWebsite := r.PostForm.Get("website")

	campaignOneTime := false
	if len(r.PostForm["oneTime"]) > 0 {
		err = json.Unmarshal([]byte(r.PostForm["oneTime"][0]), &campaignOneTime)
		if err != nil {
			c.l.Errorf("failed to unmarshal oneTime bool: %v", err)
			http.Error(w, "Failed to parse oneTime bool", http.StatusBadRequest)
			return
		}
	}
	campaignStart := r.PostForm.Get("campaignStart")
	campaignEnd := r.PostForm.Get("campaignEnd")

	resp, err := c.cc.CreatePost(context.Background(), &prcontent.CreatePostRequest{
		UserId:          profile.UserId,
		Campaign:        isCampaign,
		AgeGroup:        ageGroup,
		CampaignOneTime: campaignOneTime,
		CampaignStart:   campaignStart,
		CampaignEnd:     campaignEnd,
		CampaignWebsite: campaignWebsite,
	})
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
			UserTags:    userTags,
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

	message, err := u.uc.Follow(context.Background(), &prusers.FollowRequest{Username: profileRequest, ToFollow: dto.ProfileToFollow})
	if err != nil {
		u.l.Printf("[ERROR] following profile: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Write([]byte(message.Message))

}

func (u *Users) Unfollow(w http.ResponseWriter, r *http.Request) {
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

	message, err := u.uc.UnFollow(context.Background(), &prusers.FollowRequest{Username: profileRequest, ToFollow: dto.ProfileToFollow})
	if err != nil {
		u.l.Printf("[ERROR] following profile: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Write([]byte(message.Message))

}

func (u *Users) FollowRespond(w http.ResponseWriter, r *http.Request) {
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

	username := claims.Username

	dto := saltdata.FollowRequestDOT{}
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

	_, err = u.uc.SetFollowRequestRespond(context.Background(), &prusers.FollowRequestRespond{
		Username:        username,
		RequestUsername: dto.RequestProfile,
		Accepted:        dto.IsAccepted,
	})
	if err != nil {
		u.l.Printf("[ERROR] follow request respond: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

}

func (a *Admin) AcceptAgent(w http.ResponseWriter, r *http.Request) {
	_, err := getUserByJWS(r, a.uc)
	if err != nil {
		a.l.Errorf("failed fetching user: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.l.Errorf("failed to get email: %v\n", err)
		http.Error(w, "No email", http.StatusBadRequest)
		return
	}

	email := string(body)

	_, err = a.uc.VerifyEmail(context.Background(), &prusers.VerifyEmailRequest{Email: email})
	if err != nil {
		a.l.Errorf("failed to verify email: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = a.ac.RemoveAgentRegistration(context.Background(), &pradmin.RemoveAgentRegistrationRequest{Email: email})
	if err != nil {
		a.l.Errorf("failed to remove agent request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Accepted"))
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

func (u *Users) MuteProfile(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameByJWS(r)
	if err != nil {
		u.l.Println("faild to parse jws %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	dto := saltdata.ProfileRequestDTO{}
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

	_, err = u.uc.MuteProfile(context.Background(), &prusers.MuteProfileRequest{Logged: username, Profile: dto.Profile})
	if err != nil {
		u.l.Errorf("failed to mute profile: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (u *Users) UnmuteProfile(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameByJWS(r)
	if err != nil {
		u.l.Println("faild to parse jws %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	dto := saltdata.ProfileRequestDTO{}
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

	_, err = u.uc.UnmuteProfile(context.Background(), &prusers.UnmuteProfileRequest{Logged: username, Profile: dto.Profile})
	if err != nil {
		u.l.Errorf("failed to unmute profile: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (u *Users) BlockProfile(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameByJWS(r)
	if err != nil {
		u.l.Println("faild to parse jws %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	dto := saltdata.ProfileRequestDTO{}
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

	_, err = u.uc.BlockProfile(context.Background(), &prusers.BlockProfileRequest{Logged: username, Profile: dto.Profile})
	if err != nil {
		u.l.Errorf("failed to block profile: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (u *Users) UnblockProfile(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameByJWS(r)
	if err != nil {
		u.l.Println("faild to parse jws %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	dto := saltdata.ProfileRequestDTO{}
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

	_, err = u.uc.UnblockProfile(context.Background(), &prusers.UnblockProfileRequest{Logged: username, Profile: dto.Profile})
	if err != nil {
		u.l.Errorf("failed to unblock profile: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (u *Users) AddCloseFriend(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameByJWS(r)
	if err != nil {
		u.l.Println("faild to parse jws %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	dto := saltdata.ProfileRequestDTO{}
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

	_, err = u.uc.AddCloseFriend(context.Background(), &prusers.CloseFriendRequest{Logged: username, Profile: dto.Profile})
	if err != nil {
		u.l.Errorf("failed to add close friend: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (u *Users) RemoveCloseFriend(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameByJWS(r)
	if err != nil {
		u.l.Println("faild to parse jws %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	dto := saltdata.ProfileRequestDTO{}
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

	_, err = u.uc.RemoveCloseFriend(context.Background(), &prusers.CloseFriendRequest{Logged: username, Profile: dto.Profile})
	if err != nil {
		u.l.Errorf("failed to remove close friend: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func (c *Content) SavePost(w http.ResponseWriter, r *http.Request) {

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

	dto := saltdata.PostDTO{}
	err = saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		c.l.Errorf("failure saving post: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(dto.Id, 0, 64)
	if err != nil {
		c.l.Errorf("failed to convert id post: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	_, err = c.cc.SavePost(context.Background(), &prcontent.SavePostRequest{UserId: user.Id, PostId: id})

	if err != nil {
		c.l.Errorf("failed to save post: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func (n *Notification) NotificationSeen(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameByJWS(r)
	if err != nil {
		n.l.Println("failed to parse jws %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = n.nc.NotificationSeen(context.Background(), &prnotifications.NProfile{Username: username})
	if err != nil {
		n.l.Errorf("failed to update notifications: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
