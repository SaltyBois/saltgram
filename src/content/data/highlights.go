package data

import (
	"fmt"
	"saltgram/data"
)

type Highlight struct {
	data.Identifiable
	Name string `json:"name"`
	Stories []*Media `json:"stories" gorm:"many2many:highlight_stories;"`
	UserID uint64 `json:"userId" gorm:"type:numeric"`
}

func (db *DBConn) CreateHighlight(highlight *Highlight) error {
	return db.DB.Omit("Stories").Create(highlight).Error
}

var ErrHiglightNotFound = fmt.Errorf("higlight not found")
func (db *DBConn) AddStoriesToHighlight(highlightId uint64, ids ...uint64) error {
	highlight := Highlight{}
	res := db.DB.First(&highlight, highlightId)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrHiglightNotFound
	}

	stories, err := db.GetMediaByIds(ids...)
	if err != nil {
		return err
	}

	return db.DB.Model(&highlight).Association("Stories").Append(stories)
}