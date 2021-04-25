package servers

import (
	"context"
	"log"
	"saltgram/email/data"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Email struct {
	premail.UnimplementedEmailServer
	l  *log.Logger
	db *data.DBConn
	uc prusers.UsersClient
}

func NewEmail(l *log.Logger, db *data.DBConn, uc prusers.UsersClient) *Email {
	return &Email{l: l, db: db, uc: uc}
}

func (e *Email) ConfirmReset(ctx context.Context, r *premail.ConfirmRequest) (*premail.ConfirmResponse, error) {

	email, err := data.ConfirmPasswordReset(e.db, r.Token)
	if err != nil {
		e.l.Printf("[ERROR] confirming password reset: %v\n", err)
		return &premail.ConfirmResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &premail.ConfirmResponse{Email: email}, nil
	// TODO(Jovan): Move out to api
	// cookie := http.Cookie{
	// 	Name:     "email",
	// 	Value:    email,
	// 	Expires:  time.Now().UTC().AddDate(0, 6, 0),
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteNoneMode,
	// }
	// http.SetCookie(w, &cookie)
	// w.Write([]byte("Activated"))
}

func (e *Email) RequestReset(ctx context.Context, r *premail.ResetRequest) (*premail.ResetResponse, error) {
	err := data.SendPasswordReset(e.db, r.Email)
	if err != nil {
		e.l.Printf("[ERROR] sending email request: %v\n", err)
	}
	// NOTE(Jovan): Always return 200 OK as per OWASP guidelines
	return &premail.ResetResponse{}, nil
}

func (e *Email) Activate(ctx context.Context, r *premail.ActivateRequest) (*premail.ActivateResponse, error) {
	email, err := data.ActivateEmail(e.db, r.Token)
	if err != nil {
		e.l.Printf("[ERROR] activating email: %v", err)
		return &premail.ActivateResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	_, err = e.uc.VerifyEmail(context.Background(), &prusers.VerifyEmailRequest{Email: email})
	if err != nil {
		e.l.Printf("[ERROR] activating email: %v", err)
		return &premail.ActivateResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &premail.ActivateResponse{}, nil
}

func (e *Email) SendActivation(ctx context.Context, r *premail.SendActivationRequest) (*premail.SendActivationResponse, error) {

	err := data.SendActivation(e.db, r.Email)
	if err != nil {
		e.l.Printf("[ERROR] sending email activation: %v\n", err)
		return &premail.SendActivationResponse{}, status.Error(codes.Internal, "Internal server error")
	}

	return &premail.SendActivationResponse{}, nil
}
