package data

import (
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

func (db *DBConn) UpdateProfile(p *Profile) error {
	return db.DB.Save(&p).Error
}

func (db *DBConn) GetProfileByUsername(username string) (*Profile, error) {
	profile := Profile{}
	err := db.DB.Where("username = ?", username).First(&profile).Error
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
	err = db.DB.Find(&following, ids).Error
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
	err = db.DB.Find(&profiles, ids).Error
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
