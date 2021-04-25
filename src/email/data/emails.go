package data

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	HOST_URL     = "smtp.gmail.com"
	HOST_PORT    = "587"
	EMAIL_SENDER = "bezbednovic@gmail.com"
	PASSWORD     = "Bogsvama123"
	ADDRESS      = HOST_URL + ":" + HOST_PORT
)

type EmailRequest struct {
	Token      string `json:"token" gorm:"primaryKey" validate:"required"`
	Email      string `json:"email" validate:"required"`
	ValidUntil string `json:"validUntil" validate:"required"`
}

func NewEmailRequest(email string) *EmailRequest {
	return &EmailRequest{
		Token:      getUuid(),
		ValidUntil: time.Now().UTC().AddDate(0, 0, 1).String(),
		Email:      email,
	}
}

func ActivateEmail(db *DBConn, token string) (string, error) {
	a, err := db.Get(token)
	if err != nil {
		fmt.Println("Failed to activate for token: ", token)
		return "", err
	}

	if err := a.IsValid(db); err != nil {
		return "", err
	}
	fmt.Println("Activated email: ", a.Email)
	// err = verifyEmail(a.Email)
	// if err != nil {
	// 	return err
	// }
	db.Remove(a.Token)
	return a.Email, sendConfirmation(a.Email)
}

func ConfirmPasswordReset(db *DBConn, token string) (string, error) {
	r, err := db.Get(token)
	if err != nil {
		return "", err
	}
	fmt.Printf("Confirming request: %v\n%v\n", r.Token, token)

	if err := r.IsValid(db); err != nil {
		return "", err
	}
	db.Remove(r.Token)
	return r.Email, nil
}

func sendConfirmation(email string) error {
	subject := "Salty Boys: Account Activated"
	content := "Dear user,\nYour account has been activated!\n"
	return sendEmail(email, subject, content)
}

var ErrorActivationExpired = fmt.Errorf("verification email has expired")
var ErrorActivationInvalid = fmt.Errorf("verification email is invalid")

func (e *EmailRequest) IsValid(db *DBConn) error {
	layout := os.Getenv("TIME_LAYOUT")
	t, err := time.Parse(layout, e.ValidUntil)
	if err != nil {
		return err
	}

	if t.UTC().Before(time.Now().UTC()) {
		return ErrorActivationExpired
	}

	requests := []*EmailRequest{}
	// db.DB.Where("type == ?", erType).Find(&requests)

	db.DB.Find(&requests)
	for _, ve := range requests {
		if e.Token == ve.Token {
			return nil
		}
	}
	return ErrorActivationInvalid
}

func SendActivation(db *DBConn, email string) error {
	request := NewEmailRequest(email)
	db.Add(request)
	subject := "Salty Bois: Account Activation"
	content := fmt.Sprintf("Dear user,\nClick this link to activate your account: \n%s", generateActivationURL(request.Token))
	return sendEmail(email, subject, content)
}

func SendPasswordReset(db *DBConn, email string) error {
	request := NewEmailRequest(email)
	fmt.Printf("Adding request: %v\n", request.Token)
	db.Add(request)
	subject := "Salty Bois: Password Reset"
	content := fmt.Sprintf("Dear user,\nClick this link to reset your password:\n%s\n\n"+
		"If you did not request this, ignore this email!", generatePasswordChangeURL(request.Token))
	return sendEmail(email, subject, content)
}

func sendEmail(email, subject, content string) error {
	emailAuth := smtp.PlainAuth("", EMAIL_SENDER, PASSWORD, HOST_URL)
	signature := "\n\nSincerely Yours,\nThe Salty Bois"
	messageBytes := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n%s\n\n%s", email, subject, content, signature))
	em := []string{email}
	return smtp.SendMail(ADDRESS, emailAuth, EMAIL_SENDER, em, messageBytes)
}

func (db *DBConn) Add(req *EmailRequest) error {
	return db.DB.Create(req).Error
}

func (db *DBConn) Get(token string) (*EmailRequest, error) {
	request := EmailRequest{}
	err := db.DB.First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (db *DBConn) Remove(token string) error {
	req, err := db.Get(token)
	if err != nil {
		return err
	}
	return db.DB.Delete(req).Error
}

func getUuid() string {
	uuid, _ := uuid.NewRandom()
	uuidstring := strings.ReplaceAll(uuid.String(), "-", "")
	return uuidstring
}

func generateActivationURL(token string) string {
	return fmt.Sprintf("https://localhost:%s/email/activate/%s", os.Getenv("SALT_WEB_PORT"), token)
}

func generatePasswordChangeURL(token string) string {
	return fmt.Sprintf("https://localhost:%s/email/change/%s", os.Getenv("SALT_WEB_PORT"), token)
}
