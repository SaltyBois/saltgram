package data

import (
	"context"
	"fmt"
	"os"

	recaptcha "cloud.google.com/go/recaptchaenterprise/apiv1"
	"github.com/go-playground/validator"
	recaptchapb "google.golang.org/genproto/googleapis/cloud/recaptchaenterprise/v1"
)

type ReCaptcha struct {
	Token  string `json:"token" validate:"required"`
	Action string `json:"action" validate:"required"`
}

var ErrorInvalidToken = fmt.Errorf("invalid reCaptcha token")

func (r *ReCaptcha) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *ReCaptcha) Verify() (float32, error) {
	siteKey := os.Getenv("RECAPTCHA_SITE_KEY")
	assessmentName := "login_assessment"
	parentProject := os.Getenv("RECAPTCHA_PROJECT")

	ctx := context.Background()
	client, err := recaptcha.NewClient(ctx)
	if err != nil {
		return -1, err
	}

	event := &recaptchapb.Event{
		ExpectedAction: r.Action,
		Token:          r.Token,
		SiteKey:        siteKey,
	}

	assessment := &recaptchapb.Assessment{
		Event: event,
		Name:  assessmentName,
	}

	request := &recaptchapb.CreateAssessmentRequest{
		Assessment: assessment,
		Parent:     parentProject,
	}

	response, err := client.CreateAssessment(ctx, request)

	if err != nil {
		return -1, err
	}

	if !response.TokenProperties.Valid {
		return -1, fmt.Errorf("token was invalid because of following reasons: %v", response.TokenProperties.InvalidReason)
	} else {
		if response.Event.ExpectedAction == r.Action {
			return response.RiskAnalysis.Score, nil
		} else {
			return -1, fmt.Errorf("action attribute in reCaptcha tag does not match the action expected for scoring")
		}
	}
}
