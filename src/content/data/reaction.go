package data

import "saltgram/data"

type ReactionType string

const (
	LIKE    = "LIKE"
	DISLIKE = "DISLIKE"
)

type Reaction struct {
	data.Identifiable
	ReactionType string `json:"reactionType" validate:"required"`
	UserID       uint64 `json:"userId" gorm:"type:numeric"`
	Post         Post   `json:"post"`
	PostID       uint64 `json:"postId" gorm:"type:numeric"`
}

func (db *DBConn) AddReaction(reaction *Reaction) error {
	return db.DB.Create(reaction).Error
}

func (db *DBConn) GetReactionByPostId(id uint64) (*[]Reaction, error) {
	reaction := []Reaction{}
	err := db.DB.Where("post_id = ?", id).Find(&reaction).Error
	return &reaction, err
}

func (db *DBConn) GetReactionByUserId(id uint64) (*Reaction, error) {
	reaction := Reaction{}
	err := db.DB.Where("user_id = ?", id).First(&reaction).Error
	return &reaction, err
}

func (db *DBConn) GetReactionById(id uint64) (*Reaction, error) {
	vr := Reaction{}
	err := db.DB.Where("id = ?", id).First(&vr).Error
	return &vr, err
}

func (db *DBConn) UpdateReaction(r *Reaction) error {
	reaction := Reaction{}

	err := db.DB.First(&reaction).Error
	if err != nil {
		return err
	}

	return db.DB.Save(r).Error
}

func PutReaction(db *DBConn, rt string, id uint64) error {
	reaction, err := db.GetReactionById(id)
	if err != nil {
		return err
	}
	reaction.ReactionType = rt
	db.UpdateReaction(reaction)
	return nil
}
