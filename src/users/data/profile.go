package data

type Profile struct {
	Username  string     `json:"username" validate:"required" gorm:"primaryKey"`
	User      User       `gorm:"foreignKey:Username; references:Username" `
	Public    bool       `json:"-"`
	Taggable  bool       `json:"-"`
	Followers []*Profile `gorm:"many2many:profile_followers"`
}

type FollowRequest struct {
	ToFollow        string
	Profile         Profile `gorm:"foreignKey:ToFollow"`
	Follower        string
	FollowerProfile Profile `gorm:"foreignKey:Follower"`
}

func (db *DBConn) GetProfiles() []*Profile {
	profiles := []*Profile{}
	db.DB.Find(&profiles)
	return profiles
}

func (db *DBConn) AddProfile(p *Profile) error {
	// err := p.User.GenerateSaltAndHashedPassword()
	// if err != nil {
	// 	return err
	// }
	return db.DB.Create(p).Error
}

func GetProfileByUsername(db *DBConn, username string) (*Profile, error) {
	profile := Profile{}
	err := db.DB.First(&profile, username).Error
	return &profile, err
}
