package data

import (
	"fmt"
	"saltgram/data"
	"saltgram/protos/content/prcontent"
)

type Highlight struct {
	data.Identifiable
	Name    string   `json:"name"`
	Stories []*Media `json:"stories" gorm:"many2many:highlight_stories;"`
	UserID  uint64   `json:"userId" gorm:"type:numeric"`
}

func (db *DBConn) CreateHighlight(highlight *Highlight) error {
	return db.DB.Create(highlight).Error
}

var ErrHiglightNotFound = fmt.Errorf("higlight not found")

func (db *DBConn) GetHighlights(userId uint64) ([]*Highlight, error) {
	highlights := []*Highlight{}
	res := db.DB.Preload("Stories").Where("user_id = ?", userId).Find(&highlights)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, ErrHiglightNotFound
	}
	db.l.Infof("Media[0] Id: %v", highlights[0].Stories[0].ID)
	return highlights, nil
}

func DataToPRHighlight(d *Highlight) *prcontent.Highlight {
	stories := []*prcontent.Media{}
	for _, s := range d.Stories {
		stories = append(stories, DataToPRMedia(s))
	}

	return &prcontent.Highlight{
		Name:    d.Name,
		Stories: stories,
	}
}
