package data

import (
	"fmt"
	"saltgram/data"
)

type Verification struct {
	data.Identifiable
	UserID uint64 `json:"userId" gorm:"type:numeric"`
	URL    string `json:"url"`
}

func (db *DBConn) AddVerification(v *Verification) error {
	return db.DB.Create(v).Error
}

var ErrVerificationNotFound = fmt.Errorf("verification not found")

func (db *DBConn) GetVerificationByUserId(userId uint64) (*Verification, error) {
	v := Verification{}
	res := db.DB.Where("user_id = ?", userId).First(&v)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, ErrVerificationNotFound
	}
	return &v, nil
}
