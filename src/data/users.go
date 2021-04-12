package data

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             uint64    `json:"id"`
	FullName       string    `json:"fullName" validate:"required"`
	Email          string    `json:"email" validate:"required"`
	Username       string    `json:"username" validate:"required"`
	HashedPassword string    `json:"hashedPassword" validate:"required"`
	ReCaptcha      ReCaptcha `json:"reCaptcha" validate:"required"`

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

const (
	SALT_LENGTH = 10
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+[];',./{}:<>?")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (u *User) GenerateSaltAndHashedPassword() error {
	rand.Seed(time.Now().UnixNano())
	u.Salt = randSeq(SALT_LENGTH)
	var hns strings.Builder
	hns.WriteString(u.HashedPassword)
	hns.WriteString(u.Salt)
	bytes := []byte(hns.String())
	hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashedPassword = string(hash)
	return nil
}

func VerifyPassword(username, plainPassword string) (string, error) {
	user, _, err := findUserByUsername(username)
	if err != nil {
		return "", err
	}

	err = user.VerifyPassword(plainPassword)
	if err != nil {
		return "", err
	}
	return user.HashedPassword, nil
}

var ErrorInvalidPassword = fmt.Errorf("invalid password")

func (u *User) VerifyPassword(plainPassword string) error {
	var hns strings.Builder
	hns.WriteString(plainPassword)
	hns.WriteString(u.Salt)
	plainPasswordBytes := []byte(hns.String())
	hashedPasswordBytes := []byte(u.HashedPassword)
	err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, plainPasswordBytes)
	if err != nil {
		return ErrorInvalidPassword
	}
	return nil
}

// TODO(Jovan): For testing currently, use a DB of sorts
func GetUsers() Users {
	return userList
}

func AddUser(u *User) error {
	u.ID = getNextID()
	err := u.GenerateSaltAndHashedPassword()
	if err != nil {
		return err
	}
	userList = append(userList, u)
	return nil
}

func UpdateUser(id uint64, u *User) error {
	_, pos, err := findUserByID(id)
	if err != nil {
		return err
	}

	u.ID = id
	userList[pos] = u

	return nil
}

func GetUserByID(id uint64) (*User, error) {
	user, _, err := findUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	user, _, err := findUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

var ErrUserNotFound = fmt.Errorf("User not found")

func findUserByUsername(username string) (*User, int, error) {
	for i, u := range userList {
		if u.Username == username {
			return u, i, nil
		}
	}
	return nil, -1, ErrUserNotFound
}

func findUserByID(id uint64) (*User, int, error) {
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

var userList = []*User{}

func Seed() {
	smith := User{
		ID:             1,
		FullName:       "Mr Smith",
		Username:       "AgentSmith",
		Email:          "smith@email.com",
		HashedPassword: "smith123",

		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	}

	smith.GenerateSaltAndHashedPassword()

	userList = append(userList, &smith)
}
