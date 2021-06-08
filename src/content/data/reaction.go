package data

type ReactionType int

const (
	Like ReactionType = iota
	Dislike
)

type Reaction struct {
	ID           uint64       `json:"id"`
	ReactionType ReactionType `validate:"required"`
	Post         Post         `json:"post"`
	PostID       uint64       `json:"postId"`
}

func (db *DBConn) GetReactionByPostId(id string) (*Reaction, error) {
	reaction := Reaction{}
	err := db.DB.Where("post.id = ?", id).First(&reaction).Error
	return &reaction, err
}
