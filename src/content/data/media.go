package data

type Media struct {
	ID string `gorm:"primaryKey"`
	Filename string `json:"filename" validate:"required"`
	Tags []string `json:"tags" validate:"required"`
	Description string `json:"description" validate:"required"`
	AddedOn string `-`
	Location Location `gorm:"foreignKey:Id; references:Id" `
}

type SharedMedia struct {
	//Media []*Media
}

type Story struct {
	SharedMedia SharedMedia `validate:"required"`
	CloseFriends bool `json:"-"`
}

type Post struct {
	SharedMedia SharedMedia `validate:"required"`
	Public bool `json:"-"`
}