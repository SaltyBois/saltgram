package data

type ReactionType string

const (
	LIKE ReactionType = "LIKE"
	DISLIKE
)

type Reaction struct {
	ID           uint64       `json:"id"`
	ReactionType ReactionType `json:"reactionType" validate:"required"`
	User         User         `json:"user"`
	UserID       uint64       `json:"userId"`
	Post         Post         `json:"post"`
	PostID       uint64       `json:"postId"`
}

func (db *DBConn) AddReaction(reaction *Reaction) error {
	return db.DB.Create(reaction).Error
}

func (db *DBConn) GetReactionByPostId(id uint64) (*Reaction, error) {
	reaction := Reaction{}
	err := db.DB.Where("post_id = ?", id).First(&reaction).Error
	return &reaction, err
}

func (db *DBConn) GetReactionByUserId(id uint64) (*Reaction, error) {
	reaction := Reaction{}
	err := db.DB.Where("user_id = ?", id).First(&reaction).Error
	return &reaction, err
}