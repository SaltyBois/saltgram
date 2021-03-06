package data

import (
	"fmt"
	"math/rand"
	"saltgram/data"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	data.Identifiable
	Email          string    `json:"email" validate:"required" gorm:"unique"`
	FullName       string    `json:"fullName" validate:"required"`
	Username       string    `json:"username" validate:"required" gorm:"unique"`
	HashedPassword string    `json:"password" validate:"required"`
	ReCaptcha      ReCaptcha `json:"reCaptcha" gorm:"embedded" validate:"required"`
	Role           string    `json:"role"`

	Activated bool   `json:"-"`
	Salt      string `json:"-"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

type InfluencerRequest struct {
	data.Identifiable
	InfluencerID uint64 `gorm:"type:numeric" json:"influencerId"`
	CampaignID   uint64 `gorm:"type:numeric" json:"campaignId"`
	Website      string `json:"website"`
}

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

func (db *DBConn) AddInfluencerRequest(ir *InfluencerRequest) error {
	return db.DB.Create(ir).Error
}

func (db *DBConn) GetInfluencerRequests(influencerId uint64) (*[]InfluencerRequest, error) {
	ir := []InfluencerRequest{}
	err := db.DB.Find(&ir).Where("influencer_id = ?", influencerId).Error
	return &ir, err
}

func (db *DBConn) RemoveInfluencerRequest(influencerId, campaignId uint64) error {
	ir := InfluencerRequest{}
	db.DB.Find(&ir).Where("influencer_id = ? AND campaign_id = ?", influencerId, campaignId)
	return db.DB.Delete(&ir).Error
}

var ErrorNewPasswordSameAsOld = fmt.Errorf("new password same as old")

func ResetPassword(db *DBConn, email, password string) error {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return err
	}

	oldHashedPassword := user.HashedPassword

	user.HashedPassword = password
	err = user.GenerateSaltAndHashedPassword()
	if err != nil {
		return err
	}

	if oldHashedPassword == user.HashedPassword {
		return ErrorNewPasswordSameAsOld
	}

	err = db.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
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

func GetRole(db *DBConn, username string) (string, error) {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	return user.Role, nil
}

func IsEmailVerified(db *DBConn, username string) bool {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		return false
	}
	return user.Activated
}

func VerifyEmail(db *DBConn, email string) error {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return err
	}
	user.Activated = true
	db.UpdateUser(user)
	return nil
}

func ActivateUser(db *DBConn, userId uint64) error {
	user, err := db.GetUserById(userId)
	if err != nil {
		return err
	}
	user.Activated = true
	return db.UpdateUser(user)
}

func ChangePassword(db *DBConn, email, oldPlainPassword, newPlainPassword string) error {
	user, err := db.GetUserByUsername(email)
	if err != nil {
		return err
	}

	err = user.VerifyPassword(oldPlainPassword)
	if err != nil {
		return err
	}
	oldHashed := user.HashedPassword
	user.HashedPassword = newPlainPassword
	err = user.GenerateSaltAndHashedPassword()
	if err != nil {
		user.HashedPassword = oldHashed
		return err
	}
	db.UpdateUser(user)
	return nil
}

func IsPasswordValid(db *DBConn, username, plainPassword string) (string, error) {
	user, err := db.GetUserByUsername(username)
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

func (db *DBConn) GetUsers() []*User {
	users := []*User{}
	db.DB.Find(&users)
	return users
}

func (db *DBConn) AddUser(u *User) error {
	err := u.GenerateSaltAndHashedPassword()
	if err != nil {
		return err
	}
	return db.DB.Create(u).Error
}

func (db *DBConn) UpdateUser(u *User) error {
	user := User{}

	// NOTE(Jovan): Check if exists
	err := db.DB.First(&user).Error
	if err != nil {
		return err
	}

	return db.DB.Save(u).Error
}

func (db *DBConn) GetUserByEmail(email string) (*User, error) {
	user := User{}
	err := db.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (db *DBConn) GetUserByUsername(username string) (*User, error) {
	user := User{}
	err := db.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (db *DBConn) DeleteUser(username string) error {
	user := User{}
	return db.DB.Where("username = ?", username).Delete(&user).Error
}

func (db *DBConn) GetUserById(id uint64) (*User, error) {
	user := User{}
	err := db.DB.Where("id = ?", id).First(&user).Error
	return &user, err
}

//Moved to profile
/*func (db *DBConn) GetAllUsersByUsernameSubstring(username string) ([]User, error) {
	var users []User
	query := "%" + username + "%"
	err := db.DB.Where("username LIKE ?", query).Limit(21).Find(&users).Error
	return users, err
}*/

func Seed() {
	smith := User{
		FullName:       "Mr Smith",
		Username:       "AgentSmith",
		Email:          "smith@email.com",
		HashedPassword: "smith123",
		Role:           "user",

		Activated: true,
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	}

	smith.GenerateSaltAndHashedPassword()
}

type AccessClaims struct {
	Username       string             `json:"username"`
	Password       string             `json:"password"`
	StandardClaims jwt.StandardClaims `json:"standardClaims"`
}

type RefreshClaims struct {
	Username       string             `json:"username"`
	StandardClaims jwt.StandardClaims `json:"standardClaims"`
}

var ErrorEmptyClaims = fmt.Errorf("empty credentials")

func (uc AccessClaims) Valid() error {
	if len(uc.Username) <= 0 || len(uc.Password) <= 0 {
		return ErrorEmptyClaims
	}

	return uc.StandardClaims.Valid()
}

func (rc RefreshClaims) Valid() error {
	if len(rc.Username) <= 0 {
		return ErrorEmptyClaims
	}

	return rc.StandardClaims.Valid()
}
