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
	Token      string `json:"token" validate:"required"`
	Email      string `json:"email" validate:"required"`
	ValidUntil string `json:"validUntil" validate:"required"`
}

func NewEmailRequest(email string) *EmailRequest {
	return &EmailRequest{
		Token: getUuid(),
		ValidUntil: time.Now().UTC().AddDate(0, 0, 1).String(),
		Email: email,
	}
}

func ActivateEmail(token string) error {
	a, err := activationEmails.Find(token)
	if err != nil {
		fmt.Println("Failed to activate for token: ", token)
		return err
	}

	if err := a.IsValidFor(activationEmails); err != nil {
		return err
	}
	fmt.Println("Activated email: ",  a.Email)
	activationEmails.Remove(a.Token)
	return sendConfirmation(a.Email)
}

func ConfirmPasswordReset(token string) (string, error) {
	r, err := resetEmails.Find(token)
	if err != nil {
		return "", err
	}
	fmt.Printf("Confirming request: %v\n%v\n", r.Token, token)

	if err := r.IsValidFor(resetEmails); err != nil {
		return "", err
	}
	resetEmails.Remove(r.Token)
	return r.Email, nil
}

func sendConfirmation(email string) error {
	subject := "Salty Boys: Account Activated"
	content := "Dear user,\nYour account has been activated!\n"
	return sendEmail(email, subject, content)
}

var ErrorActivationExpired = fmt.Errorf("verification email has expired")
var ErrorActivationInvalid = fmt.Errorf("verification email is invalid")

func (e *EmailRequest) IsValidFor(requests Requests) error {
	layout := os.Getenv("TIME_LAYOUT")
	t, err := time.Parse(layout, e.ValidUntil)
	if err != nil {
		return err
	}

	if t.UTC().Before(time.Now().UTC()) {
		return ErrorActivationExpired
	}

	for _, ve := range requests {
		if e.Token == ve.Token {
			return nil
		}
	}
	return ErrorActivationInvalid
}

func SendActivation(email string) error {
	request := NewEmailRequest(email)
	activationEmails.Add(request)
	subject := "Salty Bois: Account Activation"
	content := fmt.Sprintf("Dear user,\nClick this link to activate your account: \n%s", generateActivationURL(request.Token))
	return sendEmail(email, subject, content)
}

func SendPasswordReset(email string) error {
	request := NewEmailRequest(email)
	fmt.Printf("Adding request: %v\n", request.Token)
	resetEmails.Add(request)
	subject := "Salty Bois: Password Reset"
	content := fmt.Sprintf("Dear user,\nClick this link to reset your password:\n%s\n\n" +
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

type Requests []*EmailRequest

var resetEmails = Requests{
	// NOTE(Jovan): pls don't hack me :(
	// {
	// 	Token: "632309ff4c154735991b58589964b7ed",
	// 	ValidUntil: time.Now().AddDate(1, 0, 0).UTC().String(),
	// 	Email: "ivos.jovan@protonmail.ch",
	// },
}
var activationEmails = Requests{
	{
		Token:      "48d135c6d7194cf08fb3ce4a31f06618",//getUuid(),
		ValidUntil: time.Now().AddDate(0, 0, 1).UTC().String(),
		Email:      "admin@email.com",
	},
}

// TODO(Jovan): Remove!
func GetAllActivations() []*EmailRequest {
	return activationEmails
}

// TODO(Jovan): Remove!
func GetAllResets() []*EmailRequest {
	return resetEmails
}

var ErrorEmailRequestNotFound = fmt.Errorf("email request not found")

func (r *Requests) Add(req *EmailRequest) {
	*r = append(*r, req)
}

func (r Requests) Find(token string) (*EmailRequest, error) {
	for _, req := range r {
		if req.Token == token {
			return req, nil
		}
	}
	return nil, ErrorEmailRequestNotFound
}

func (r *Requests) Remove(token string) {
	for i, req := range *r {
		if req.Token == token {
			(*r)[len(*r)-1], (*r)[i] = (*r)[i], (*r)[len(*r)-1]
			*r = (*r)[:len(activationEmails)-1]
		}
	}
}

// TODO(Jovan): Temp for "seeding"
func getUuid() string {
	uuid, _ := uuid.NewRandom()
	uuidstring := strings.ReplaceAll(uuid.String(), "-", "")
	return uuidstring
}

func generateActivationURL(token string) string {
	return fmt.Sprintf("http://localhost%s/email/activate/%s", os.Getenv("PORT_FRONT_SALT"), token)
}

func generatePasswordChangeURL(token string) string {
	return fmt.Sprintf("http://localhost%s/email/change/%s", os.Getenv("PORT_FRONT_SALT"), token)
}
