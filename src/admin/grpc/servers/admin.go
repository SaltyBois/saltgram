package servers

import (
	"context"
	"saltgram/admin/data"
	"saltgram/protos/admin/pradmin"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Admin struct {
	pradmin.UnimplementedAdminServer
	l  *logrus.Logger
	db *data.DBConn
}

func NewAdmin(l *logrus.Logger, db *data.DBConn) *Admin {
	return &Admin{
		l:  l,
		db: db,
	}
}

func (a *Admin) GetPendingVerifications(ctx context.Context, r *pradmin.GetVerificationRequest) (*pradmin.GetVerificationResponse, error) {
	verificationRequests, err := a.db.GetPendingVerificationRequests()
	if err != nil {
		a.l.Errorf("failure getting verification requests: %v\n", err)
		return &pradmin.GetVerificationResponse{}, err
	}
	requests := []*pradmin.VerificationRequest{}
	for _, vr := range *verificationRequests {
		requests = append(requests, &pradmin.VerificationRequest{
			Id:       vr.ID,
			FullName: vr.Fullname,
			Category: vr.Category,
			Url:      vr.URL,
			UserId:   vr.UserID,
		})
	}
	return &pradmin.GetVerificationResponse{VerificationRequest: requests}, nil
}

func (a *Admin) AddVerificationReq(ctx context.Context, r *pradmin.AddVerificationRequest) (*pradmin.AddVerificationResponse, error) {

	verificationRequest := data.VerificationRequest{
		Fullname: r.FullName,
		UserID:   r.UserId,
		Category: r.Category,
		Status:   data.PENDING,
		URL:      r.Url,
	}
	err := a.db.AddVerificationRequest(&verificationRequest)
	if err != nil {
		a.l.Errorf("failure adding verification request: %v\n", err)
		return &pradmin.AddVerificationResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &pradmin.AddVerificationResponse{}, nil
}

func (a *Admin) ReviewVerificationReq(ctx context.Context, r *pradmin.ReviewVerificatonRequest) (*pradmin.ReviewVerificatonResponse, error) {

	err := data.ReviewVerificationRequest(a.db, r.Status, r.Id)
	if err != nil {
		a.l.Errorf("failure updating verification request: %v\n", err)
		return &pradmin.ReviewVerificatonResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &pradmin.ReviewVerificatonResponse{}, nil
}

func (a *Admin) SendInappropriateContentReport(ctx context.Context, r *pradmin.InappropriateContentReportRequest) (*pradmin.InappropriateContentReportResponse, error) {

	inappropriateContentReport := data.InappropriateContentReport{
		UserID:        r.UserId,
		Status:        data.PENDING,
		SharedMediaID: r.SharedMediaId,
	}
	err := a.db.AddInappropriateContentReport(&inappropriateContentReport)
	if err != nil {
		a.l.Errorf("failure sending inappropriate content report: %v\n", err)
		return &pradmin.InappropriateContentReportResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &pradmin.InappropriateContentReportResponse{}, nil
}

func (a *Admin) GetPendingInappropriateContentReport(ctx context.Context, r *pradmin.GetInappropriateContentReportRequest) (*pradmin.GetInappropriateContentReportResponse, error) {
	reports, err := a.db.GetPendingInappropriateContentReport()
	if err != nil {
		a.l.Errorf("failure getting inappropriate report: %v\n", err)
		return &pradmin.GetInappropriateContentReportResponse{}, err
	}
	reps := []*pradmin.InappropriateContentReport{}
	for _, vr := range *reports {
		reps = append(reps, &pradmin.InappropriateContentReport{
			Id:            vr.ID,
			UserId:        vr.UserID,
			SharedMediaId: vr.SharedMediaID,
		})
	}
	return &pradmin.GetInappropriateContentReportResponse{InappropriateContentReport: reps}, nil
}
