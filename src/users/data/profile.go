package data

import (
	"saltgram/data"
	"time"
)

type RequestStatus string

const (
	PENDING RequestStatus = "PENDING"
	ACCEPTED
	DENIED
)

type Profile struct {
	data.Identifiable
	UserID          uint64
	User            User
	Username        string     `json:"username" gorm:"unique"`
	Public          bool       `json:"isPublic"`
	Taggable        bool       `json:"isTaggable"`
	Description     string     `json:"description"`
	Following       []*Profile `gorm:"many2many:profile_following;"`
	Profiles        []FollowRequest
	Requests        []FollowRequest
	PhoneNumber     string    `json:"phoneNumber"`
	Gender          string    `json:"gender"`
	DateOfBirth     time.Time `json:"dateOfBirth"`
	WebSite         string    `json:"webSite"`
	PrivateProfile  bool      `json:"privateProfile"` // Why
	ProfileFolderId string    `json:"-"`
	PostsFolderId   string    `json:"-"`
	StoriesFolderId string    `json:"-"`
	ProfilePictureURL string `json:"profilePictureURL"`
}

type FollowRequest struct {
	ID            uint64        `json:"-"`
	ProfileID     uint64        `json:"profileId"`
	RequestID     uint64        `json:"followerId"`
	RequestStatus RequestStatus `json:"stats"`
}

//TODO(Marko add profliPicture?)
type ProfileFollowerDTO struct {
	Username string
	FullName string
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

func (db *DBConn) UpdateProfilePicture(url, username string) (error) {
	profile, err := db.GetProfileByUsername(username)
	if err != nil {
		return err
	}
	profile.ProfilePictureURL = url
	return db.UpdateProfile(profile)
}

func (db *DBConn) UpdateProfile(p *Profile) error {
	return db.DB.Save(&p).Error
}

func (db *DBConn) GetProfileByUsername(username string) (*Profile, error) {
	profile := Profile{}
	err := db.DB.Where("username = ?", username).First(&profile).Error
	return &profile, err
}

func CheckIfFollowing(db *DBConn, profile_username string, following_user_id uint64) (bool, error) {
	var count int64
	// err := db.DB.Table("profile_followers").Where("profile_username = ? AND follower_username = ?", profile, username).Count(&count).Error
	err := db.DB.Raw("SELECT * FROM profile_following LEFT JOIN profiles on profile_id = user_id WHERE username = ? AND following_id = ? ", profile_username, following_user_id).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func GetFollowerCount(db *DBConn, username string) (int64, error) {
	var count int64
	// err := db.DB.Table("profile_followers").Where("follower_username = ?", username).Count(&count).Error
	err := db.DB.Raw("SELECT * FROM profile_following LEFT JOIN profiles on following_id = user_id WHERE username = ?", username).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetFollowingCount(db *DBConn, username string) (int64, error) {
	var count int64
	// err := db.DB.Table("profile_followers").Where("profile_username = ?", username).Count(&count).Error
	err := db.DB.Raw("SELECT * FROM profile_following LEFT JOIN profiles on profile_id = user_id WHERE username = ?", username).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func SetFollow(db *DBConn, profile *Profile, profileToFollow *Profile) error {
	db.DB.Model(&profile).Association("Following").Append(&profileToFollow)
	return db.DB.Save(&profile).Error
}

func CreateFollowRequest(db *DBConn, profile *Profile, request *Profile) error {
	db.DB.Model(&profile).Association("Requests").Append(&request)
	return db.DB.Save(&request).Error
}

func GetFollowers(db *DBConn, username string) ([]Profile, error) {

	var followers []Profile
	//err := db.DB.Preload("Followers").Where("follower_username = ?", username).Find(&followers).Error
	err := db.DB.Raw("SELECT * FROM profile_following LEFT JOIN profiles on following_id = user_id WHERE username = ?", username).Scan(&followers).Error
	if err != nil {
		return nil, err
	}
	return followers, nil
}

func GetFollowing(db *DBConn, username string) ([]Profile, error) {
	var following []Profile
	err := db.DB.Raw("SELECT * FROM profile_following LEFT JOIN profiles on profile_id = user_id WHERE username = ?", username).Scan(&following).Error
	if err != nil {
		return nil, err
	}
	return following, nil
}
