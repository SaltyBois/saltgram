package data

type Profile struct {
	Username    string     `json:"username" validate:"required" gorm:"primaryKey"`
	User        User       `gorm:"foreignKey:Username; references:Username" `
	Public      bool       `json:"isPublic"`
	Taggable    bool       `json:"-"`
	Description string     `json:"description"`
	Followers   []*Profile `gorm:"many2many:profile_followers"`
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
	err := db.DB.Where("username = ?", username).First(&profile).Error
	return &profile, err
}

func CheckIfFollowing(db *DBConn, profile string, username string) (bool, error) {
	var count int64
	err := db.DB.Table("profile_followers").Where("profile_username = ? AND follower_username = ?", profile, username).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func GetFollowerCount(db *DBConn, username string) (int64, error) {
	var count int64
	err := db.DB.Table("profile_followers").Where("follower_username = ?", username).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetFollowingCount(db *DBConn, username string) (int64, error) {
	var count int64
	err := db.DB.Table("profile_followers").Where("profile_username = ?", username).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func SetFollow(db *DBConn, profile *Profile, profileToFollow *Profile) error {
	db.DB.Model(&profileToFollow).Association("Followers").Append(&profile)
	return db.DB.Save(&profileToFollow).Error
}
