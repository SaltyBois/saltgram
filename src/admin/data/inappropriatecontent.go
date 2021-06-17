package data

import "saltgram/data"

type InappropriateContentReport struct {
	data.Identifiable
	SharedMediaID uint64 `json:"sharedMediaId"`
	UserID        uint64 `json:"userId"`
	Status        string `json:"status"`
}

func (db *DBConn) AddInappropriateContentReport(inappropriateContentReport *InappropriateContentReport) error {
	return db.DB.Create(inappropriateContentReport).Error
}
