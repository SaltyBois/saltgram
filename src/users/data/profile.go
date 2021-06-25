package data

import (
	"fmt"
	"saltgram/data"
	"time"
)

type RequestStatus string

const (
	PENDING  RequestStatus = "PENDING"
	ACCEPTED RequestStatus = "ACCEPTED"
	DENIED   RequestStatus = "DENIED"
)

type Profile struct {
	data.Identifiable
	UserID            uint64
	User              User
	Username          string     `json:"username" gorm:"unique"`
	Public            bool       `json:"isPublic"`
	Taggable          bool       `json:"isTaggable"`
	Messagable        bool       `json:"messageable"`
	Description       string     `json:"description"`
	Following         []*Profile `gorm:"many2many:profile_following;"`
	Profiles          []FollowRequest
	Requests          []FollowRequest
	PhoneNumber       string    `json:"phoneNumber"`
	Gender            string    `json:"gender"`
	DateOfBirth       time.Time `json:"dateOfBirth"`
	WebSite           string    `json:"webSite"`
	PrivateProfile    bool      `json:"privateProfile"` // Why
	ProfileFolderId   string    `json:"-"`
	PostsFolderId     string    `json:"-"`
	StoriesFolderId   string    `json:"-"`
	ProfilePictureURL string    `json:"profilePictureURL"`
	AccountType       string    `json:"accountType"`
	Verified          bool      `json:"verified"`
	Active            bool      `json:"-"`
}

type FollowRequest struct {
	data.Identifiable
	ProfileID     uint64        `json:"profileId" gorm:"type:numeric"`
	RequestID     uint64        `json:"followerId" gorm:"type:numeric"`
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

func (db *DBConn) VerifyProfile(userId uint64, accountType string) error {
	p, err := db.GetProfileByUserId(userId)
	if err != nil {
		return err
	}
	p.Verified = true
	p.AccountType = accountType
	return db.UpdateProfile(p)
}

func (db *DBConn) UpdateProfilePicture(url, username string) error {
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

var ErrProfileNotFound = fmt.Errorf("profile not found")

func (db *DBConn) GetProfileByUserId(userId uint64) (*Profile, error) {
	profile := Profile{}
	res := db.DB.Where("user_id = ? AND active = ?", userId, true).First(&profile)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, ErrProfileNotFound
	}

	return &profile, nil
}

func (db *DBConn) GetProfileByUsername(username string) (*Profile, error) {
	profile := Profile{}
	err := db.DB.Where("username = ? AND active = ?", username, true).First(&profile).Error
	return &profile, err
}

func CheckIfFollowing(db *DBConn, profile *Profile, following *Profile) (bool, error) {
	var count int64
	err := db.DB.Table("profile_following").Where("profile_id = ? AND following_id = ?", profile.ID, following.ID).Count(&count).Error
	//err := db.DB.Raw("SELECT COUNT(*) FROM profile_following LEFT JOIN profiles on profile_user_id = user_id WHERE username = ? AND following_user_id = ? ", profile_username, following_user_id).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func GetFollowerCount(db *DBConn, profile *Profile) (int64, error) {
	var count int64
	err := db.DB.Table("profile_following").Where("following_id = ?", profile.ID).Count(&count).Error
	//err := db.DB.Raw("SELECT COUNT(*) FROM profile_following LEFT JOIN profiles on following_user_id = user_id WHERE username = ?", username).Count(&count).Error
	if err != nil {
		return 0, err
	}
	// if count == 0 {
	// 	return -1, nil
	// }
	return count, nil
}

func GetFollowingCount(db *DBConn, profile *Profile) (int64, error) {
	var count int64
	err := db.DB.Table("profile_following").Where("profile_id = ?", profile.ID).Count(&count).Error
	//err := db.DB.Raw("SELECT COUNT(*) FROM profile_following LEFT JOIN profiles on profile_id = user_id WHERE username = ?", username).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func SetFollow(db *DBConn, profile *Profile, profileToFollow *Profile) error {
	return db.DB.Exec("INSERT INTO profile_following (profile_id, following_id) VALUES (?, ?)", profile.ID, profileToFollow.ID).Error
}

func Unfollow(db *DBConn, profile *Profile, profileToUnfollow *Profile) error {
	return db.DB.Exec("DELETE FROM profile_following WHERE profile_id = ? AND following_id = ?", profile.ID, profileToUnfollow.ID).Error
}

func CreateFollowRequest(db *DBConn, profile *Profile, request *Profile) error {
	//profile.Requests = append(profile.Requests, request)
	fr := FollowRequest{
		ProfileID:     profile.ID,
		RequestID:     request.ID,
		RequestStatus: PENDING,
	}
	return db.DB.Create(&fr).Error
}

func GetFollowers(db *DBConn, profile *Profile) ([]Profile, error) {
	var followers []Profile
	var ids []uint64
	err := db.DB.Table("profile_following").Select("profile_id").Where("following_id = ?", profile.ID).Find(&ids).Error
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return followers, nil
	}
	err = db.DB.Find(&followers, ids).Error

	if err != nil {
		return nil, err
	}
	return followers, nil
}

func GetFollowing(db *DBConn, profile *Profile) ([]Profile, error) {
	var following []Profile
	var ids []uint64
	err := db.DB.Table("profile_following").Select("following_id").Where("profile_id = ?", profile.ID).Find(&ids).Error
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return following, nil
	}
	//TODO(Marko) Check later
	err = db.DB.Where("active = ?", true).Find(&following, ids).Error
	if err != nil {
		return nil, err
	}
	return following, nil
}

func GetFollowRequests(db *DBConn, profile *Profile) ([]Profile, error) {
	var profiles []Profile
	var ids []uint64
	err := db.DB.Table("follow_requests").Select("request_id").Where("profile_id = ? AND request_status = ?", profile.ID, PENDING).Find(&ids).Error
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return profiles, nil
	}
	//TODO(Marko) Check where later
	err = db.DB.Where("active = ?", true).Find(&profiles, ids).Error
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func FollowRequestRespond(db *DBConn, profile *Profile, profile_request *Profile, accepted bool) error {
	fr := FollowRequest{}
	err := db.DB.Where("profile_id = ? AND request_id = ? AND request_status = ?", profile.ID, profile_request.ID, PENDING).First(&fr).Error
	if err != nil {
		return err
	}
	if accepted {
		fr.RequestStatus = ACCEPTED
	} else {
		fr.RequestStatus = DENIED
	}
	return db.DB.Save(&fr).Error
}

func CheckForFollowingRequest(db *DBConn, profile *Profile, profile_request *Profile) (bool, error) {
	var count int64
	err := db.DB.Model(&FollowRequest{}).Where("profile_id = ? AND request_id = ? AND request_status = ?", profile.ID, profile_request.ID, PENDING).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func (db *DBConn) GetAllProfilesByUsernameSubstring(username string) ([]Profile, error) {
	var profiles []Profile
	query := "%" + username + "%"
	err := db.DB.Where("username LIKE ?", query).Limit(21).Find(&profiles).Error
	return profiles, err
}
