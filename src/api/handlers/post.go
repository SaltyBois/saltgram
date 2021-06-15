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
		Username: dto.Username,
		FullName: dto.FullName,
		Email:    dto.Email,
		Password: dto.Password,
		ReCaptcha: &prusers.UserReCaptcha{
			Token:  dto.ReCaptcha.Token,
			Action: dto.ReCaptcha.Action,
		},
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

func (c *Content) AddSharedMedia(w http.ResponseWriter, r *http.Request) {

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

	dto := saltdata.SharedMediaDTO{}
	err = saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		c.l.Errorf("failure adding shared media: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	media := []*prcontent.Media{}
	for _, m := range dto.Media {
		tags := []*prcontent.Tag{}
		for _, t := range m.Tags {
			tags = append(tags, &prcontent.Tag{
				Value: t.Value,
				Id:    t.ID,
			})
		}
		media = append(media, &prcontent.Media{
			UserId:      user.Id,
			Filename:    m.Filename,
			Tags:        tags,
			Description: m.Description,
			Location: &prcontent.Location{
				Country: m.Location.Country,
				State:   m.Location.State,
				ZipCode: m.Location.ZipCode,
				Street:  m.Location.Street,
			},
			AddedOn: m.AddedOn,
		})
	}

	_, err = c.cc.AddSharedMedia(context.Background(), &prcontent.AddSharedMediaRequest{Media: media})

	if err != nil {
		c.l.Errorf("failed to add shared media: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func (a *Admin) AddVerificationRequest(w http.ResponseWriter, r *http.Request) {

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

	dto := saltdata.VerificationRequestDTO{}
	err = saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		a.l.Errorf("failure adding verification request: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = a.ac.AddVerificationReq(context.Background(), &pradmin.AddVerificationRequest{FullName: dto.FullName, UserId: user.Id, Category: dto.Category})

	if err != nil {
		a.l.Errorf("failed to add verification request: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}
