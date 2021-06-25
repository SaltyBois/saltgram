package servers

import (
	"context"
	"saltgram/admin/data"
	"saltgram/protos/admin/pradmin"
	"saltgram/protos/content/prcontent"
	"saltgram/protos/users/prusers"
	"strconv"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Admin struct {
	pradmin.UnimplementedAdminServer
	l  *logrus.Logger
	db *data.DBConn
	cc prcontent.ContentClient
	uc prusers.UsersClient
}

func NewAdmin(l *logrus.Logger, db *data.DBConn, cc prcontent.ContentClient, uc prusers.UsersClient) *Admin {
	return &Admin{
		l:  l,
		db: db,
		cc: cc,
		uc: uc,
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

	userId, category, err := data.ReviewVerificationRequest(a.db, r.Status, r.Id)
	if err != nil {
		a.l.Errorf("failure updating verification request: %v\n", err)
		return &pradmin.ReviewVerificatonResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	if r.Status == "ACCEPTED" {
		_, err = a.uc.VerifyProfile(ctx, &prusers.VerifyProfileRequest{UserId: userId, AccountType: category})
		if err != nil {
			a.l.Errorf("failed to update profile: %v", err)
			return &pradmin.ReviewVerificatonResponse{}, status.Error(codes.Internal, "Internal error")
		}
	}

	return &pradmin.ReviewVerificatonResponse{}, nil
}

func (a *Admin) SendInappropriateContentReport(ctx context.Context, r *pradmin.InappropriateContentReportRequest) (*pradmin.InappropriateContentReportResponse, error) {

	resp, err := a.cc.GetPostPreviewURL(ctx, &prcontent.GetPostPreviewURLRequest{PostId: r.PostId})
	if err != nil {
		a.l.Errorf("failed to get shared media url: %v", err)
		return &pradmin.InappropriateContentReportResponse{}, err
	}

	inappropriateContentReport := data.InappropriateContentReport{
		UserID: r.UserId,
		Status: data.PENDING,
		PostID: r.PostId,
		URL:    resp.Url,
	}

	err = a.db.AddInappropriateContentReport(&inappropriateContentReport)
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
			Id:     strconv.FormatUint(vr.ID, 10),
			UserId: vr.UserID,
			PostId: strconv.FormatUint(vr.PostID, 10),
			Url:    vr.URL,
		})
	}
	return &pradmin.GetInappropriateContentReportResponse{InappropriateContentReport: reps}, nil
}

func (a *Admin) RejectInappropriateContentReport(ctx context.Context, r *pradmin.RejectInappropriateContentReportRequest) (*pradmin.RejectInappropriateContentReportResponse, error) {

	i, err := strconv.ParseUint(r.Id, 10, 64)
	if err != nil {
		a.l.Errorf("failure parsing id: %v\n", err)
		return &pradmin.RejectInappropriateContentReportResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	err = data.RejectInappropriateContentReport(a.db, i)
	if err != nil {
		a.l.Errorf("failure rejecting inappropriate content report: %v\n", err)
		return &pradmin.RejectInappropriateContentReportResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &pradmin.RejectInappropriateContentReportResponse{}, nil
}

func (a *Admin) AcceptInappropriateContentReport(ctx context.Context, r *pradmin.AcceptInappropriateContentReportRequest) (*pradmin.AcceptInappropriateContentReportResponse, error) {

	i, err := strconv.ParseUint(r.Id, 10, 64)
	if err != nil {
		a.l.Errorf("failure parsing id: %v\n", err)
		return &pradmin.AcceptInappropriateContentReportResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	err = data.AcceptInappropriateContentReport(a.db, i)
	if err != nil {
		a.l.Errorf("failure accepting inappropriate content report: %v\n", err)
		return &pradmin.AcceptInappropriateContentReportResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &pradmin.AcceptInappropriateContentReportResponse{}, nil
}
