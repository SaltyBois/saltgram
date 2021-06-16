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

func (a *Admin) GetPendingVerifications(r *pradmin.GetVerificationRequest, stream pradmin.Admin_GetPendingVerificationsServer) error {
	verificationRequests, err := a.db.GetPendingVerificationRequests()
	if err != nil {
		a.l.Errorf("failure getting verification requests: %v\n", err)
		return err
	}
	for _, vr := range *verificationRequests {
		err = stream.Send(&pradmin.GetVerificationResponse{
			Id:       vr.ID,
			FullName: vr.Fullname,
			Category: vr.Category,
			//Media: vr.Media,
		})
		if err != nil {
			a.l.Errorf("failure sending verification request response: %v\n", err)
			return err
		}
	}
	return nil
}

func (a *Admin) AddVerificationReq(ctx context.Context, r *pradmin.AddVerificationRequest) (*pradmin.AddVerificationResponse, error) {

	verificationRequest := data.VerificationRequest{
		Fullname: r.FullName,
		UserID:   r.UserId,
		Category: r.Category,
		Status:   data.PENDING,
		//Media:    r.Media,
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
		UserID: r.UserId,
		Status: data.PENDING,
		Reason: r.Reason,
		//SharedMedia:    r.SharedMedia,
	}
	err := a.db.AddInappropriateContentReport(&inappropriateContentReport)
	if err != nil {
		a.l.Errorf("failure sending inappropriate content report: %v\n", err)
		return &pradmin.InappropriateContentReportResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &pradmin.InappropriateContentReportResponse{}, nil
}
