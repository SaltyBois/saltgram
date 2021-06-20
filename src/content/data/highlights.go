package data

import (
	"fmt"
	"saltgram/data"
)

type Highlight struct {
	data.Identifiable
	Name string `json:"name"`
	Stories []*Media `json:"stories"`
	UserID uint64 `json:"userId" gorm:"type:numeric"`
}

func (db *DBConn) CreateHighlight(highlight *Highlight) error {
	return db.DB.Create(highlight).Error
}

var ErrHiglightNotFound = fmt.Errorf("higlight not found")
func (db *DBConn) AddStoriesToHighlight(highlightId uint64, stories []*Media) error {
	highlight := Highlight{}
	res := db.DB.First(&highlight, highlightId)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrHiglightNotFound
	}

	highlight.Stories = append(highlight.Stories, stories...)
	return db.DB.Save(&highlight).Error
}