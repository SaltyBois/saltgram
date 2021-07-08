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
	REJECTED = "REJECTED"
	PENDING  = "PENDING"
)

type VerificationRequest struct {
	data.Identifiable
	Fullname string `json:"fullname" validate:"required"`
	URL      string `json:"url"`
	UserID   uint64 `json:"userId" gorm:"type:numeric"`
	Category string `json:"category" validate:"required"`
	Status   string `json:"status" validate:"required"`
}

type AgentRegistrationRequest struct {
	data.Identifiable
	AgentEmail string `json:"agentEmail"`
}

func (db *DBConn) GetAgentRegistrations() (*[]AgentRegistrationRequest, error) {
	ar := []AgentRegistrationRequest{}
	err := db.DB.Find(&ar).Error
	return &ar, err
}

func (db *DBConn) AddAgentRegistrationRequest(ar *AgentRegistrationRequest) error {
	return db.DB.Create(ar).Error
}

func (db *DBConn) RemoveAgentRegistrationRequest(email string) error {
	ar := AgentRegistrationRequest{}
	db.DB.First(&ar).Where("agent_email = ?", email)
	return db.DB.Delete(&ar).Error
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

func ReviewVerificationRequest(db *DBConn, status string, id uint64) (uint64, string, error) {
	verificationRequest, err := db.GetVerficationRequestById(id)
	if err != nil {
		return 0, "", err
	}
	verificationRequest.Status = status
	db.UpdateVerificationRequest(verificationRequest)
	return verificationRequest.UserID, verificationRequest.Category, nil
}
