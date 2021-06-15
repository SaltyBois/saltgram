package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	saltdata "saltgram/data"
	"saltgram/protos/auth/prauth"
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
		e.l.Printf("[ERROR] getting email: %v\n", err)
		http.Error(w, "No email", http.StatusBadRequest)
		return
	}
	email := string(body)
	go func() {
		_, err = e.ec.RequestReset(context.Background(), &premail.ResetRequest{Email: email})
		if err != nil {
			e.l.Printf("[ERROR] sending email request: %v\n", err)
		}
	}()
	// NOTE(Jovan): Always send 200 OK as per OWASP
	w.Write([]byte("200 - OK"))
}

func (u *Users) ChangePassword(w http.ResponseWriter, r *http.Request) {
	cr := ChangeRequest{}

	err := saltdata.FromJSON(&cr, r.Body)
	if err != nil {
		u.l.Printf("[ERROR] deserializing ChangeRequest: %v\n", err)
		http.Error(w, "Failed to parse request", http.StatusBadRequest)
		return
	}

	err = cr.Validate()
	if err != nil {
		u.l.Printf("[ERROR] ChangeRequest not valid: %v\n", err)
		http.Error(w, "Bad change request", http.StatusBadRequest)
		return
	}
	jws, err := getUserJWS(r)
	if err != nil {
		u.l.Printf("[ERROR] getting jws: %v\n", err)
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

	_, err = u.uc.ChangePassword(context.Background(), &prusers.ChangeRequest{
		Username:         claims.Username,
		OldPlainPassword: cr.OldPassword,
		NewPlainPassword: cr.NewPassword,
	})

	if err != nil {
		u.l.Printf("[ERROR] POST change password request: %v\n", err)
		http.Error(w, "Error in POST change password request", http.StatusBadRequest)
		return
	}
	w.Write([]byte("200 - OK"))
}

func (u *Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		u.l.Printf("[ERROR] parsing body: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("emailforreset")
	if err != nil {
		u.l.Printf("[ERROR] getting cookie: %v", err)
		http.Error(w, "No reset request cookie", http.StatusBadRequest)
		return
	}

	newPassword := string(body)
	_, err = u.uc.ResetPassword(context.Background(), &prusers.UserResetRequest{Email: cookie.Value, Password: newPassword})
	if err != nil {
		u.l.Printf("[ERROR] resetting password: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Write([]byte("200 - OK"))
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	dto := saltdata.UserDTO{}
	err := saltdata.FromJSON(&dto, r.Body)
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

	_, err = u.uc.Register(context.Background(), &prusers.RegisterRequest{
		Username: dto.Username,
		FullName: dto.FullName,
		Email:    dto.Email,
		Password: dto.Password,
		Description: dto.Description,
		ReCaptcha: &prusers.UserReCaptcha{
			Token:  dto.ReCaptcha.Token,
			Action: dto.ReCaptcha.Action,
		},
		PhoneNumber: dto.PhoneNumber,
		Gender: dto.Gender,
		DateOfBirth: dto.DateOfBirth.Unix(),
		WebSite: dto.WebSite,
		PrivateProfile: dto.PrivateProfile,
	})

	if err != nil {
		u.l.Printf("[ERROR] registering user: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Activation email sent"))
}

func (a *Auth) GetJWT(w http.ResponseWriter, r *http.Request) {
	user := saltdata.Login{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		a.l.Printf("[ERROR] deserializing user: %v", err)
		http.Error(w, "Failed to deserialize user", http.StatusBadRequest)
		return
	}
	res, err := a.ac.GetJWT(context.Background(), &prauth.JWTRequest{Username: user.Username, Password: user.Password})
	if err != nil {
		a.l.Printf("[ERROR] getting jwt: %v\n", err)
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
		a.l.Printf("[ERROR] deserializing body: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = login.Validate()
	if err != nil {
		a.l.Printf("[ERROR] validating: %v\n", err)
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
		a.l.Printf("[ERROR] calling login: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	saltdata.ToJSON(res, w)
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

func (u *Users) UpdateProfile(w http.ResponseWriter, r *http.Request) {
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

	dto := saltdata.ProflieDTO{}
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

	_, err = u.uc.UpdateProfile(context.Background(), &prusers.UpdateRequest{
		OldUsername: username,
		NewUsername: dto.Username,
		Email:       dto.Email,
		FullName:    dto.FullName,
		Public:      dto.Public,
		Taggable:    dto.Taggable,
	})
	if err != nil {
		u.l.Printf("[ERROR] updating profile: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Updated %s"))

}
