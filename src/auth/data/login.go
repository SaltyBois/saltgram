package data

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
)

// NOTE(Jovan): Refresh token
type Refresh struct {
	Username string `json:"username" gorm:"primaryKey" validate:"required"`
	Token    string `json:"token" validate:"required"`
}

type RefreshClaims struct {
	Username       string             `json:"username"`
	StandardClaims jwt.StandardClaims `json:"standardClaims"`
}

func (r *Refresh) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

var ErrorEmptyClaims = fmt.Errorf("empty credentials")

func (rc RefreshClaims) Valid() error {
	if len(rc.Username) <= 0 {
		return ErrorEmptyClaims
	}

	return rc.StandardClaims.Valid()
}

func AddRefreshToken(db *DBConn, username, token string) error {
	return db.DB.Create(&Refresh{username, token}).Error
}

func UpdateRefreshTokenUsername(db *DBConn, oldUsername string, newUsername string) error {
	return db.DB.Model(&Refresh{}).Where("username = ?", oldUsername).Update("username", newUsername).Error
}

var ErrorRefreshTokenNotFound = fmt.Errorf("refresh token not found")

func GetRefreshToken(db *DBConn, username string) (string, error) {
	r := Refresh{}
	err := db.DB.First(&r).Where("username = ?", username).Error
	if err != nil {
		return "", ErrorRefreshTokenNotFound
	}
	return r.Token, nil
}

func GetRefreshTokens(db *DBConn) []*Refresh {
	tokens := []*Refresh{}
	db.DB.Find(&tokens)
	return tokens
}

func (r *Refresh) Verify(db *DBConn) error {
	// NOTE(Jovan): https://gorm.io/docs/security.html
	rt := Refresh{}
	return db.DB.Where("TOKEN = ?", r.Token).First(&rt).Error
}

func DeleteRefreshToken(db *DBConn, username string) error {
	r := Refresh{}
	return db.DB.Where("username = ?", username).Delete(&r).Error
}
