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

	"github.com/gorilla/mux"
)

func (a *Auth) CheckPermissions(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.l.Printf("[ERROR] reading body: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	route := string(body)
	jws, _ := getUserJWS(r)
	_, err = a.ac.CheckPermissions(context.Background(),
		&prauth.PermissionRequest{
			Jws:    jws,
			Path:   route,
			Method: r.Method,
		})
	if err != nil {
		a.l.Printf("[ERROR] authenticating: %v\n", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	w.Write([]byte("200 - OK"))
}

func (e *Email) ConfirmReset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	res, err := e.ec.ConfirmReset(context.Background(), &premail.ConfirmRequest{Token: token})
	if err != nil {
		e.l.Errorf("failure confirming reset: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	cookie := http.Cookie{
		Name:     "emailforreset",
		Value:    res.Email,
		Expires:  time.Now().UTC().AddDate(0, 6, 0),
		HttpOnly: true,
		Secure:   true,
		Path:     "/users",
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("200 - OK"))
}

func (e *Email) Activate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	_, err := e.ec.Activate(context.Background(), &premail.ActivateRequest{Token: token})
	if err != nil {
		e.l.Errorf("failure activating email: %v\n", err)
		http.Error(w, "Failed to activate email", http.StatusBadRequest)
		return
	}

	w.Write([]byte("200 - OK"))
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
		OldUsername:    username,
		NewUsername:    dto.Username,
		Email:          dto.Email,
		FullName:       dto.FullName,
		Public:         !dto.PrivateProfile,
		Taggable:       dto.Taggable,
		Description:    dto.Description,
		PhoneNumber:    dto.PhoneNumber,
		Gender:         dto.Gender,
		DateOfBirth:    dto.DateOfBirth.Unix(),
		WebSite:        dto.WebSite,
		PrivateProfile: dto.PrivateProfile,
		Messageable:    dto.Messageable,
	})

	if err != nil {
		u.l.Printf("[ERROR] updating profile: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Updated %s"))

}

func (a *Auth) UpdateJWTUsername(w http.ResponseWriter, r *http.Request) {

	// TODO: PART 1 JWS
	jws, err := getUserJWS(r)
	if err != nil {
		a.l.Println("[ERROR] JWS not found")
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
		a.l.Printf("[ERROR] parsing claims: %v", err)
		http.Error(w, "Error parsing claims", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*saltdata.AccessClaims)

	if !ok {
		a.l.Println("[ERROR] unable to parse claims")
		http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
		return
	}

	username := claims.Username
	password := claims.Password

	type NewUsernameDto struct {
		NewUsername string `json:"newUsername" validate:"required"`
	}

	// TODO: PART 2 UPDATING

	dto := NewUsernameDto{}
	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		a.l.Printf("[ERROR] deserializing JWTDto: %v", err)
		http.Error(w, "Failed to deserialize JWTDto", http.StatusBadRequest)
		return
	}

	_, err = a.ac.UpdateRefresh(context.Background(), &prauth.UpdateRefreshRequest{OldUsername: username, NewUsername: dto.NewUsername})
	if err != nil {
		a.l.Printf("[ERROR] updating jwt: %v\n", err)
		http.Error(w, "Failed to update jwt", http.StatusBadRequest)
		return
	}

	// TODO: PART 3 SETTING NEW UPDATED TOKEN

	res, err := a.ac.GetJWT(context.Background(), &prauth.JWTRequest{Username: dto.NewUsername, Password: password})
	if err != nil {
		a.l.Printf("[ERROR] getting new updated jwt: %v\n", err)
		http.Error(w, "Failed to get new updated jwt", http.StatusBadRequest)
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

func (c *Content) PutReaction(w http.ResponseWriter, r *http.Request) {

	dto := saltdata.ReactionPutDTO{}
	err := saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		c.l.Errorf("failure updating: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = c.cc.PutReaction(context.Background(), &prcontent.PutReactionRequest{Id: dto.Id, ReactionType: dto.ReactionType})

	if err != nil {
		c.l.Errorf("failed to updating: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func (a *Admin) ReviewVerificationRequest(w http.ResponseWriter, r *http.Request) {

	dto := saltdata.ReviewRequestDTO{}
	err := saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		a.l.Errorf("failure reviewing verification request: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	i, err := strconv.ParseUint(dto.Id, 10, 64)

	if err != nil {
		a.l.Errorf("failed to convert id string: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = a.ac.ReviewVerificationReq(context.Background(), &pradmin.ReviewVerificatonRequest{Id: i, Status: dto.Status})

	if err != nil {
		a.l.Errorf("failed to review verification request: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func (a *Admin) RejectInappropriateContentReport(w http.ResponseWriter, r *http.Request) {

	dto := saltdata.ReviewReportDTO{}
	err := saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		a.l.Errorf("failure rejecting inappropriate content report: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = a.ac.RejectInappropriateContentReport(context.Background(), &pradmin.RejectInappropriateContentReportRequest{Id: dto.Id})

	if err != nil {
		a.l.Errorf("failed to reject inappropriate content report: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func (a *Admin) RemoveInappropriateContent(w http.ResponseWriter, r *http.Request) {

	dto := saltdata.ReviewReportDTO{}
	err := saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		a.l.Errorf("failure rejecting inappropriate content report: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	i, err := strconv.ParseUint(dto.SharedMediaId, 10, 64)

	if err != nil {
		a.l.Errorf("failed to convert id string: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = a.cc.DeleteSharedMedia(context.Background(), &prcontent.DeleteSharedMediaRequest{Id: i})

	if err != nil {
		a.l.Errorf("failed to reject inappropriate content report: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = a.ac.AcceptInappropriateContentReport(context.Background(), &pradmin.AcceptInappropriateContentReportRequest{Id: dto.Id})

	if err != nil {
		a.l.Errorf("failed to reject inappropriate content report: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

}

func (a *Admin) RemoveProfile(w http.ResponseWriter, r *http.Request) {

	dto := saltdata.ReviewReportDTO{}
	err := saltdata.FromJSON(&dto, r.Body)
	if err != nil {
		a.l.Errorf("failure rejecting inappropriate content report: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = a.ac.AcceptInappropriateContentReport(context.Background(), &pradmin.AcceptInappropriateContentReportRequest{Id: dto.Id})

	if err != nil {
		a.l.Errorf("failed to reject inappropriate content report: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	resp, err := a.cc.GetPostUserId(context.Background(), &prcontent.GetPostUserIdRequest{PostId: dto.SharedMediaId})

	if err != nil {
		a.l.Errorf("failed to get user id by post: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := a.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: resp.UserId})

	if err != nil {
		a.l.Errorf("failed fetching user %v\n", err)
		http.Error(w, "User getting error", http.StatusInternalServerError)
		return
	}

	/*profile, err := a.uc.GetProfileByUsername(context.Background(), &prusers.ProfileRequest{User: user.Username, Username: user.Username})
	if err != nil {
		a.l.Errorf("failed fetching profile %v\n", err)
		http.Error(w, "Profile getting error", http.StatusInternalServerError)
		return
	}*/

	_, err = a.uc.DeleteProfile(context.Background(), &prusers.Profile{Username: user.Username})
	if err != nil {
		a.l.Errorf("failed deleting profile %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
