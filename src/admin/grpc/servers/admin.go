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
		//Media:    r.Media,
		Category: r.Category,
		Status:   data.ACCEPTED,
	}
	err := a.db.AddVerificationRequest(&verificationRequest)
	if err != nil {
		a.l.Errorf("failure adding verification request: %v\n", err)
		return &pradmin.AddVerificationResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &pradmin.AddVerificationResponse{}, nil
}

/*func (a *Admin) ReviewVerificationReq(ctx context.Context, r *pradmin.ReviewVerificatonRequest) (*pradmin.ReviewVerificatonResponse, error) {

	err := ReviewVerificationRequest(a.db, r.Status, r.Id)
	if err != nil {
		a.l.Errorf("failure updating verification request: %v\n", err)
		return &pradmin.ReviewVerificatonResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &pradmin.ReviewVerificatonResponse{}, nil
}*/
