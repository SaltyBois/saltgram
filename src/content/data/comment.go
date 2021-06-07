package data

type Comment struct {
	ID string `json: "id" gorm:"primaryKey" validate:"required"`
	Content string `json:"content" validate:"required"`
	Likes int64 `json:"likes" validate:"required"`
	Dislikes int64 `json:"dislikes" validate:"required"`
	Comments []*Comment `gorm:"many2many:comment_replies"`
	Post Post `validate:"required"`
}

func (db *DBConn) Add(comment *Comment) error {
	return db.DB.Create(comment).Error
}

func (db *DBConn) GetCommentByPostId(id string) (*Comment, error) {
	comment := Comment{}
	err := db.DB.Where("post.id = ?", id).First(&comment).Error
	return &comment, err
}

