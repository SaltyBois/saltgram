package data

type ReactionType int

const (
	Like ReactionType = iota
	Dislike
)

type Reaction struct {
	ReactionType ReactionType `validate:"required"`
	//Profile Profile
	Post Post
}

func (db *DBConn) GetReactionByPostId(id string) (*Reaction, error) {
	reaction := Reaction{}
	err := db.DB.Where("post.id = ?", id).First(&reaction).Error
	return &reaction, err
}
