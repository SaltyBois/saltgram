package data

import (
	"time"

	"github.com/go-playground/validator"
)

type UserDTO struct {
	Id             string    `json:"id"`
	Email          string    `json:"email" validate:"required"`
	FullName       string    `json:"fullName" validate:"required"`
	Username       string    `json:"username" validate:"required"`
	Password       string    `json:"password" validate:"required"`
	Description    string    `json:"description" validate:"required"`
	ReCaptcha      ReCaptcha `json:"reCaptcha" validate:"required"`
	PhoneNumber    string    `json:"phoneNumber" validate:"required"`
	Gender         string    `json:"gender" validate:"required"`
	DateOfBirth    time.Time `json:"dateOfBirth" validate:"required"`
	WebSite        string    `json:"webSite"`
	PrivateProfile bool      `json:"privateProfile"`
	ProfilePictureURL string `json:"profilePictureURL"`
}

type ProflieDTO struct {
	Email          string    `json:"email" validate:"required"`
	FullName       string    `json:"fullName" validate:"required"`
	Username       string    `json:"username" validate:"required"`
	Public         bool      `json:"public"`      /* `json:"public" validate:"required"` */
	Taggable       bool      `json:"taggable"`    /* `json:"taggable" validate:"required"` */
	Messageable    bool      `json:"messageable"` /* `json:"taggable" validate:"required"` */
	Description    string    `json:"description" validate:"required"`
	PhoneNumber    string    `json:"phoneNumber" validate:"required"`
	Gender         string    `json:"gender" validate:"required"`
	DateOfBirth    time.Time `json:"dateOfBirth" validate:"required"`
	WebSite        string    `json:"webSite"`
	PrivateProfile bool      `json:"privateProfile"`
	UserId         string    `json:"userId"`
	Followers      int64     `json:"followers"`
	Following      int64     `json:"following"`
	IsPublic       bool      `json:"isPublic"`
	IsFollowing    bool      `json:"isFollowing"`
}

type ProfileDTO struct {
	Email             string `json:"email" validate:"required"`
	FullName          string `json:"fullName" validate:"required"`
	Username          string `json:"username" validate:"required"`
	Public            bool   `json:"public"`      /* `json:"public" validate:"required"` */
	Taggable          bool   `json:"taggable"`    /* `json:"taggable" validate:"required"` */
	Messageable       bool   `json:"messageable"` /* `json:"taggable" validate:"required"` */
	Description       string `json:"description" validate:"required"`
	PhoneNumber       string `json:"phoneNumber" validate:"required"`
	Gender            string `json:"gender" validate:"required"`
	DateOfBirth       int64  `json:"dateOfBirth" validate:"required"`
	WebSite           string `json:"webSite"`
	PrivateProfile    bool   `json:"privateProfile"`
	UserId            string `json:"userId"`
	Followers         int64  `json:"followers"`
	Following         int64  `json:"following"`
	IsPublic          bool   `json:"isPublic"`
	IsFollowing       bool   `json:"isFollowing"`
	ProfilePictureURL string `json:"profilePictureURL"`
	AccountType       string `json:"accountType"`
	Verified          bool   `json:"verified"`
	IsThisMe          bool   `json:"isThisMe"`
}

type FollowRequestDOT struct {
	RequestProfile string `json:"profile" validate:"required"`
	IsAccepted     bool   `json:"accepted"`
}

type ProfileFollowDetailedDTO struct {
	Username       string 		`json:"username"`
	Following      bool   		`json:"following"`
	Pending        bool   		`json:"pending"`
	ProfliePicture string 		`json:"profilePictureURL"`
	Id 			   string 		`json:"id"`
}

type FollowDTO struct {
	ProfileToFollow string `json:"profile" validate:"required"`
}

func (p *ProflieDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func (u *UserDTO) Validate() error {
	// TODO(Jovan): Extract into a global validator?
	validate := validator.New()
	return validate.Struct(u)
}

func (f *FollowDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(f)
}

func (fr *FollowRequestDOT) Validate() error {
	validate := validator.New()
	return validate.Struct(fr)
}
