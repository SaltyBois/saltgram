package data

type Profile struct {
	Username string `json:"username validate:required"`
	User     User   `gorm:"foreignKey:Username; references:Username" `
	Public   bool   `json:"-"`
	Taggable bool   `json:"-"`
}

func (db *DBConn) GetProfiles() []*Profile {
	profiles := []*Profile{}
	db.DB.Find(&profiles)
	return profiles
}

func (db *DBConn) AddProfile(p *Profile) error {
	return db.DB.Create(p).Error
}

func (db *DBConn) GetProfileByUsername(username string) (*Profile, error) {
	profile := Profile{}
	err := db.DB.First(&profile, username).Error
	return &profile, err
}
