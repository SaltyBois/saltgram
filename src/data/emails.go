package data

import (
	"fmt"
	"net/smtp"
	"os"
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

type Activation struct {
	Token      string `json:"token" validate:"required"`
	ValidUntil string `json:"validUntil"`
}

func (e *Activation) Activate() error {
	if err := e.IsValid(); err != nil {
		return err
	}
	removeActivation(e.Token)
	return nil
}

var ErrorActivationExpired = fmt.Errorf("verification email has expired")
var ErrorActivationInvalid = fmt.Errorf("verification email is invalid")


func (e *Activation) IsValid() error {
	layout := "2006-01-02 15:04:05 -0700 MST"
	t, err := time.Parse(layout, e.ValidUntil)
	if err != nil {
		return err
	}

	if t.UTC().Before(time.Now().UTC()) {
		return ErrorActivationExpired
	}

	for _, ve := range emails {
		if e.Token == ve.Token {
			return nil
		}
	}
	return ErrorActivationInvalid
}

func SendActivation(email string) error {
	ver := Activation{
		Token:      getUuid(),
		ValidUntil: time.Now().UTC().String(),
	}
	emails = append(emails, &ver)

	emailAuth := smtp.PlainAuth("", EMAIL_SENDER, PASSWORD, HOST_URL)
	message := []byte(
		"To: " + email + "\r\n" +
			"Subject: " + "Salty Boys: Account Activation" + "\r\n" +
			"Dear user,\nClick this link to activate your account: \n" + getActivationURL(ver.Token) +
			"\n\n Sincerely Yours,\nMe & The Salty Boyz")

	em := []string{email}

	err := smtp.SendMail(ADDRESS, emailAuth, EMAIL_SENDER, em, message)
	return err
}

var emails = []*Activation{
	{
		Token:      getUuid(),
		ValidUntil: time.Now().AddDate(0, 0, 1).UTC().String(),
	},
}

type Activations []*Activation

func removeActivation(token string) {
	for i, a := range emails {
		if a.Token == token {
			emails[len(emails)-1], emails[i] = emails[i], emails[len(emails)-1]
			emails = emails[:len(emails)-1]
		}
	}
}

// TODO(Jovan): Temp for "seeding"
func getUuid() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}

func getActivationURL(token string) string {
	return fmt.Sprintf("http://localhost%s/activate/%s", os.Getenv("PORT_SALT"), token)
}
