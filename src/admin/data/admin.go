package data

type Category string

const (
	INFLUENCER Category = "INFLUENCER"
	SPORTS
	MEDIA
	BUSINESS
	BRAND
	ORGANIZATION
)

type VerificationStatus string

const (
	ACCEPTED VerificationStatus = "ACCEPTED"
	REJECTED
	PENDING
)

type Admin struct {
}

type VerificationRequest struct {
	ID       uint64 `json:"id"`
	Fullname string `json:"fullname" validate:"required"`
	//DocumentImage image
	User               User               `json:"user"`
	UserID             uint64             `json:"userId"`
	Category           Category           `validate:"required"`
	VerificationStatus VerificationStatus `validate:"required"`
}

func (db *DBConn) AddVerificationRequest(verificationRequest *VerificationRequest) error {
	return db.DB.Create(verificationRequest).Error
}

func (db *DBConn) GetPendingVerificationRequests() (*[]VerificationRequest, error) {
	verificationRequest := []VerificationRequest{}
	err := db.DB.Where("verification_status = 'PENDING'").Find(&verificationRequest).Error
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

func ReviewVerificationRequest(db *DBConn, verificationStatus VerificationStatus, id uint64) error {
	verificationRequest, err := db.GetVerficationRequestById(id)
	if err != nil {
		return err
	}
	verificationRequest.VerificationStatus = verificationStatus
	db.UpdateVerificationRequest(verificationRequest)
	return nil
}
