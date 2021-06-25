package data

import "saltgram/data"

type InappropriateContentReport struct {
	data.Identifiable
	PostID uint64 `json:"sharedMediaId" gorm:"type:numeric"`
	UserID uint64 `json:"userId" gorm:"type:numeric"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

func (db *DBConn) AddInappropriateContentReport(inappropriateContentReport *InappropriateContentReport) error {
	return db.DB.Create(inappropriateContentReport).Error
}

func (db *DBConn) GetPendingInappropriateContentReport() (*[]InappropriateContentReport, error) {
	inappropriateContentReport := []InappropriateContentReport{}
	err := db.DB.Where("status = 'PENDING'").Find(&inappropriateContentReport).Error
	return &inappropriateContentReport, err
}

func (db *DBConn) UpdateInappropriateContentReport(i *InappropriateContentReport) error {
	report := InappropriateContentReport{}

	err := db.DB.First(&report).Error
	if err != nil {
		return err
	}

	return db.DB.Save(i).Error
}

func (db *DBConn) GetInappropriateContentReportById(id uint64) (*InappropriateContentReport, error) {
	i := InappropriateContentReport{}
	err := db.DB.Where("id = ?", id).First(&i).Error
	return &i, err
}

func RejectInappropriateContentReport(db *DBConn, id uint64) error {
	report, err := db.GetInappropriateContentReportById(id)
	if err != nil {
		return err
	}
	report.Status = "REJECTED"
	db.UpdateInappropriateContentReport(report)
	return nil
}

func AcceptInappropriateContentReport(db *DBConn, id uint64) error {
	report, err := db.GetInappropriateContentReportById(id)
	if err != nil {
		return err
	}
	report.Status = "ACCEPTED"
	db.UpdateInappropriateContentReport(report)
	return nil
}
