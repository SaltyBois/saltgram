package data

import "saltgram/data"

type Comment struct {
	data.Identifiable
	Content  string `json:"content" validate:"required"`
	Likes    int64  `json:"likes" validate:"required"`
	Dislikes int64  `json:"dislikes" validate:"required"`
	UserID   uint64 `json:"userId"`
	PostID   uint64 `json:"postId"`
	Post     Post   `json:"post" validate:"required"`
}

func (db *DBConn) AddComment(comment *Comment) error {
	return db.DB.Create(comment).Error
}

func (db *DBConn) GetCommentByPostId(id uint64) (*Comment, error) {
	comment := Comment{}
	err := db.DB.Where("post_id = ?", id).First(&comment).Error
	return &comment, err
}
