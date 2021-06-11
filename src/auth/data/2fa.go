package data

import (
	"encoding/base32"
	"net/url"

	"github.com/dgryski/dgoogauth"
	"rsc.io/qr"
)

var (
	secret = []byte{'H', 'e', 'l', 'l', 'o', '!', 0xDE, 0xAD, 0xBE, 0xEF}
	secretB32 = base32.StdEncoding.EncodeToString(secret)
	issuer = "saltgram"
)

func Get2FAQR(email string) ([]byte, error) {
	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		return nil, err
	}
	URL.Path += "/" + url.PathEscape(issuer) + ":" + url.PathEscape(email)

	params := url.Values{}
	params.Add("secret", secretB32)
	params.Add("issuer", issuer)

	URL.RawQuery = params.Encode()
	code, err := qr.Encode(URL.String(), qr.Q)
	if err != nil {
		return nil, err
	}

	b := code.PNG()
	return b, nil
}

func Authenticate2FA(token string) (bool, error) {
	otpc := &dgoogauth.OTPConfig {
		Secret: secretB32,
		WindowSize: 3,
		HotpCounter: 0,
	}

	ok, err := otpc.Authenticate(token)
	if err != nil {
		return false, err
	}
	return ok, err
}