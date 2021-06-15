package data

type InappropriateContentReport struct {
	ID            uint        `json:"id"`
	SharedMedia   SharedMedia `json:"sharedMedia"`
	SharedMediaID uint64      `json:"sharedMediaId"`
	User          User        `json:"user"`
	UserID        uint64      `json:"userId"`
	Status        string      `json:"status"`
	Reason        string      `json:"reason"`
}

func (db *DBConn) AddInappropriateContentReport(inappropriateContentReport *InappropriateContentReport) error {
	return db.DB.Create(inappropriateContentReport).Error
}
