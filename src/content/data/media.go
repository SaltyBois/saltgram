package data

import udata "saltgram/users/data"

type Media struct {
	ID            uint64   `json:"id"`
	SharedMediaID uint64   `json:"sharedMediaId"`
	Filename      string   `json:"filename" validate:"required"`
	Tags          []Tag    `gorm:"many2many:media_tags" json:"tags" validate:"required"`
	Description   string   `json:"description" validate:"required"`
	AddedOn       string   `json:"addedOn"`
	Location      Location `gorm:"embedded"`
}

type Tag struct {
	ID    uint64 `json:"id"`
	Value string `json:"value" validate:"required"`
}

type SharedMedia struct {
	ID    uint64   `json:"id"`
	Media []*Media `json:"media"`
}

type Story struct {
	ID            uint64      `json:"id"`
	User          udata.User  `json:"user"`
	UserID        string      `json:"userId"`
	SharedMedia   SharedMedia `json:"sharedMedia"`
	SharedMediaID uint64      `json:"sharedMediaId"`
	CloseFriends  bool        `json:"closeFriends"`
}

type Post struct {
	ID            uint64      `json:"id"`
	User          udata.User  `json:"user"`
	UserID        string      `json:"userId"`
	SharedMedia   SharedMedia `validate:"required"`
	SharedMediaID uint64      `json:"sharedMediaId"`
	Public        bool        `json:"-"`
}

func (db *DBConn) GetSharedMediaByProfile(id string) (*SharedMedia, error) {
	sharedMedia := SharedMedia{}
	err := db.DB.Where("user_id = ?", id).First(&sharedMedia).Error
	return &sharedMedia, err
}

func (db *DBConn) AddSharedMedia(s *SharedMedia) error {
	return db.DB.Create(SharedMedia{}).Error
}
