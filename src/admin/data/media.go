package data

type SharedMedia struct {
	Media []*Media `json:"media"`
}

type Media struct {
	ID            uint64 `json:"id"`
	SharedMediaID uint64 `json:"sharedMediaId"`
	Filename      string `json:"filename" validate:"required"`
}
