package data

import (
	"fmt"

	"github.com/go-playground/validator"
)

// NOTE(Jovan): Refresh token
type Refresh struct {
	Username string `json:"username" gorm:"primaryKey" validate:"required"`
	Token    string `json:"token" validate:"required"`
}

func (r *Refresh) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func AddRefreshToken(db *DBConn, username, token string) error {
	return db.DB.Create(&Refresh{username, token}).Error
}

var ErrorRefreshTokenNotFound = fmt.Errorf("refresh token not found")

func GetRefreshToken(db *DBConn, username string) (string, error) {
	r := Refresh{}
	err := db.DB.First(&r).Error
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
	return db.DB.Where("TOKEN == ?", r.Token).First(&rt).Error
}
