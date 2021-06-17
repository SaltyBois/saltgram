package data

type InappropriateContentReport struct {
	ID            uint        `json:"id"`
	SharedMedia   SharedMedia `json:"sharedMedia"`
	SharedMediaID uint64      `json:"sharedMediaId"`
	UserID        uint64      `json:"userId"`
	Status        string      `json:"status"`
}

func (db *DBConn) AddInappropriateContentReport(inappropriateContentReport *InappropriateContentReport) error {
	return db.DB.Create(inappropriateContentReport).Error
}
