package data

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	SALT_LENGTH = 10
)

type User struct {
	Identifiable
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Salt     string `json:"-"`
	Agent    bool   `json:"agent"`
	Token    string `json:"token"`
}

type AccessClaims struct {
	Username       string             `json:"username"`
	Password       string             `json:"password"`
	StandardClaims jwt.StandardClaims `json:"standardClaims"`
}

var ErrorEmptyClaims = fmt.Errorf("empty credentials")

func (uc AccessClaims) Valid() error {
	if len(uc.Username) <= 0 || len(uc.Password) <= 0 {
		return ErrorEmptyClaims
	}

	return uc.StandardClaims.Valid()
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+[];',./{}:<>?")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (db *DBConn) AddUser(u *User) error {
	err := u.GenerateSaltAndHashedPassword()
	if err != nil {
		return err
	}
	return db.DB.Create(u).Error
}

func (db *DBConn) GetByUsername(username string) (*User, error) {
	u := User{}
	err := db.DB.First(&u).Where("username = ?", username).Error
	return &u, err
}

func (db *DBConn) Login(username, password string) (string, error) {
	u, err := db.GetByUsername(username)
	if err != nil {
		return "", err
	}
	err = u.VerifyPassword(password)
	if err != nil {
		return "", err
	}

	claims := AccessClaims{
		Username: username,
		Password: u.Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Minute * 30).Unix(),
			Issuer:    "Agent",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jws, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return jws, nil
}

func (db *DBConn) GetByJWS(jws string) (*User, error) {
	token, err := jwt.ParseWithClaims(
		jws,
		AccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed parsing claims: %v", err)
	}

	claims, ok := token.Claims.(*AccessClaims)

	if !ok {
		return nil, fmt.Errorf("unable to parse claims")
	}
	return db.GetByUsername(claims.Username)
}

var ErrorInvalidPassword = fmt.Errorf("invalid password")

func (u *User) VerifyPassword(plainPassword string) error {
	var hns strings.Builder
	hns.WriteString(plainPassword)
	hns.WriteString(u.Salt)
	plainPasswordBytes := []byte(hns.String())
	hashedPasswordBytes := []byte(u.Password)
	err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, plainPasswordBytes)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GenerateSaltAndHashedPassword() error {
	rand.Seed(time.Now().UnixNano())
	u.Salt = randSeq(SALT_LENGTH)
	var hns strings.Builder
	hns.WriteString(u.Password)
	hns.WriteString(u.Salt)
	bytes := []byte(hns.String())
	hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}
