package data

import udata "saltgram/users/data"

type Comment struct {
	ID       uint64     `json:"id" validate:"required"`
	Content  string     `json:"content" validate:"required"`
	Likes    int64      `json:"likes" validate:"required"`
	Dislikes int64      `json:"dislikes" validate:"required"`
	User     udata.User `json:"user"`
	UserID   string     `json:"userId"`
	PostID   uint64     `json:"postId"`
	Post     Post       `json:"post" validate:"required"`
}

func (db *DBConn) Add(comment *Comment) error {
	return db.DB.Create(comment).Error
}

func (db *DBConn) GetCommentByPostId(id uint64) (*Comment, error) {
	comment := Comment{}
	err := db.DB.Where("postId = ?", id).First(&comment).Error
	return &comment, err
}
