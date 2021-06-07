package data

type ReactionType int

const (
	Like ReactionType = iota
	Dislike
)

type Reaction struct {
	ReactionType ReactionType `validate:"required"`
}