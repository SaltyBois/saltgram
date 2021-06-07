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
	//Profile Profile
	Media []*Media
}

type Story struct {
	SharedMedia SharedMedia `validate:"required"`
	CloseFriends bool `json:"-"`
}

type Post struct {
	SharedMedia SharedMedia `validate:"required"`
	Public bool `json:"-"`
}

/*func (db *DBConn) GetSharedMediaByUsername(username string) (*SharedMedia, error) {
	sharedMedia := SharedMedia{}
	err := db.DB.Where("username = ?", username).First(&sharedMedia).Error
	return &sharedMedia, err
}*/

func (db *DBConn) GetSharedMediaByProfile(username string) (*SharedMedia, error) {
	sharedMedia := SharedMedia{}
	err := db.DB.Where("username = ?", username).First(&sharedMedia).Error
	return &sharedMedia, err
}

func (db *DBConn) AddSharedMedia(s *SharedMedia) error {
	return db.DB.Create(SharedMedia{}).Error
}