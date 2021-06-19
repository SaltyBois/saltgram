package data

import "saltgram/data"

type Category string

const (
	INFLUENCER Category = "INFLUENCER"
	SPORTS
	MEDIA
	BUSINESS
	BRAND
	ORGANIZATION
)

type Status string

const (
	ACCEPTED = "ACCEPTED"
	REJECTED
	PENDING
)

type VerificationRequest struct {
	data.Identifiable
	Fullname string `json:"fullname" validate:"required"`
	URL      string `json:"url"`
	UserID   uint64 `json:"userId"`
	Category string `json:"category" validate:"required"`
	Status   string `json:"status" validate:"required"`
}

func (db *DBConn) AddVerificationRequest(verificationRequest *VerificationRequest) error {
	return db.DB.Create(verificationRequest).Error
}

func (db *DBConn) GetPendingVerificationRequests() (*[]VerificationRequest, error) {
	verificationRequest := []VerificationRequest{}
	err := db.DB.Where("status = 'PENDING'").Find(&verificationRequest).Error
	return &verificationRequest, err
}

func (db *DBConn) UpdateVerificationRequest(v *VerificationRequest) error {
	verificationRequest := VerificationRequest{}

	err := db.DB.First(&verificationRequest).Error
	if err != nil {
		return err
	}

	return db.DB.Save(v).Error
}

func (db *DBConn) GetVerficationRequestById(id uint64) (*VerificationRequest, error) {
	vr := VerificationRequest{}
	err := db.DB.Where("id = ?", id).First(&vr).Error
	return &vr, err
}

func ReviewVerificationRequest(db *DBConn, status string, id uint64) error {
	verificationRequest, err := db.GetVerficationRequestById(id)
	if err != nil {
		return err
	}
	verificationRequest.Status = status
	db.UpdateVerificationRequest(verificationRequest)
	return nil
}
