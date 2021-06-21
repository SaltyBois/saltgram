package data

import "saltgram/data"

type InappropriateContentReport struct {
	data.Identifiable
	SharedMediaID uint64 `json:"sharedMediaId" gorm:"type:numeric"`
	UserID        uint64 `json:"userId" gorm:"type:numeric"`
	Status        string `json:"status"`
}

func (db *DBConn) AddInappropriateContentReport(inappropriateContentReport *InappropriateContentReport) error {
	return db.DB.Create(inappropriateContentReport).Error
}

func (db *DBConn) GetPendingInappropriateContentReport() (*[]InappropriateContentReport, error) {
	inappropriateContentReport := []InappropriateContentReport{}
	err := db.DB.Where("status = 'PENDING'").Find(&inappropriateContentReport).Error
	return &inappropriateContentReport, err
}
