package data

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
	User          User        `json:"user"`
	UserID        uint64      `json:"userId"`
	SharedMedia   SharedMedia `json:"sharedMedia"`
	SharedMediaID uint64      `json:"sharedMediaId"`
	CloseFriends  bool        `json:"closeFriends"`
}

type Post struct {
	ID            uint64      `json:"id"`
	User          User        `json:"user"`
	UserID        uint64      `json:"userId"`
	SharedMedia   SharedMedia `validate:"required"`
	SharedMediaID uint64      `json:"sharedMediaId"`
}

type ProfilePicture struct {
	ID      uint64 `json:"id"`
	User    User   `json:"user"`
	UserID  uint64 `json:"userId"`
	Media   Media  `validate:"required"`
	MediaID uint64 `json:"mediaId"`
	URL     string `json:"url"`
}

func (db *DBConn) GetSharedMediaByUser(id uint64) (*[]SharedMedia, error) {
	sharedMedia := []SharedMedia{}
	err := db.DB.Where("user_id = ?", id).Find(&sharedMedia).Error
	return &sharedMedia, err
}

func (db *DBConn) AddSharedMedia(s *SharedMedia) error {
	return db.DB.Create(s).Error
}

func (db *DBConn) GetStoryByUser(id uint64) (*[]Story, error) {
	story := []Story{}
	err := db.DB.Where("user_id = ?", id).Find(&story).Error
	return &story, err
}

func (db *DBConn) AddStory(s *Story) error {
	return db.DB.Create(s).Error
}
func (db *DBConn) GetPostByUser(id uint64) (*[]Post, error) {
	post := []Post{}
	err := db.DB.Where("user_id = ?", id).Find(&post).Error
	return &post, err
}

func (db *DBConn) AddPost(p *Post) error {
	return db.DB.Create(p).Error
}

func (db *DBConn) GetProfilePictureByUser(id uint64) (*ProfilePicture, error) {
	post := ProfilePicture{}
	err := db.DB.Where("user_id = ?", id).Find(&post).Error
	return &post, err
}

func (db *DBConn) AddProfilePicture(pp *ProfilePicture) error {
	oldPP, err := db.GetProfilePictureByUser(pp.UserID)
	if err != nil {
		return db.DB.Create(pp).Error
	}
	oldPP.Media = pp.Media
	oldPP.MediaID = pp.MediaID
	oldPP.URL = pp.URL
	return db.DB.Save(&oldPP).Error
}

func (db *DBConn) GetPostsByReaction(userId uint64) (*[]Post, error) {
	post := []Post{}
	err := db.DB.Raw("SELECT p.* FROM posts p INNER JOIN reactions r on p.id = r.post_id WHERE r.user_id = ?", userId).Find(&post).Error
	return &post, err
}
