package data

import (
	udata "saltgram/users/data"
)

type ReactionType int

const (
	Like ReactionType = iota
	Dislike
)

type Reaction struct {
	ID           uint64       `json:"id"`
	ReactionType ReactionType `validate:"required"`
	User         udata.User   `json:"user"`
	UserID       uint64       `json:"userId"`
	Post         Post         `json:"post"`
	PostID       uint64       `json:"postId"`
}

func (db *DBConn) GetReactionByPostId(id uint64) (*Reaction, error) {
	reaction := Reaction{}
	err := db.DB.Where("post_id = ?", id).First(&reaction).Error
	return &reaction, err
}
