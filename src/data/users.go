package data

import (
	"fmt"
	"time"

	"github.com/go-playground/validator"
)

type User struct {
	ID        uint64    `json:"id"`
	FullName  string    `json:"fullName" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	ReCaptcha ReCaptcha `json:"reCaptcha" validate:"required"`

	Salt      string `json:"-"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

// NOTE(Jovan): Collection of users
type Users []*User

func (u *User) Validate() error {
	// TODO(Jovan): Extract into a global validator?
	validate := validator.New()
	return validate.Struct(u)
}

// TODO(Jovan): For testing currently, use a DB of sorts
func GetUsers() Users {
	return userList
}

func AddUser(u *User) {
	u.ID = getNextID()
	userList = append(userList, u)
}

func UpdateUser(id uint64, u *User) error {
	_, pos, err := findUser(id)
	if err != nil {
		return err
	}

	u.ID = id
	userList[pos] = u

	return nil
}

func GetUserByID(id uint64) (*User, error) {
	user, _, err := findUser(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

var ErrUserNotFound = fmt.Errorf("User not found")

func findUser(id uint64) (*User, int, error) {
	for i, u := range userList {
		if u.ID == id {
			return u, i, nil
		}
	}

	return nil, -1, ErrUserNotFound
}

func getNextID() uint64 {
	u := userList[len(userList)-1]
	return u.ID + 1
}

var userList = []*User{
	{
		ID:        1,
		Username:  "AgentSmith",
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	},
	{
		ID:        2,
		Username:  "Neo",
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	},
}
