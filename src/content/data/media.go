package data

import (
	"fmt"
	"saltgram/data"
	"saltgram/protos/content/prcontent"
	"strconv"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EMimeType int32

const (
	EMimeType_IMAGE = iota
	EMimeType_VIDEO
)

type Media struct {
	data.Identifiable
	SharedMediaID uint64    `json:"sharedMediaId" gorm:"type:numeric"`
	Filename      string    `json:"filename" validate:"required"`
	Tags          []Tag     `gorm:"many2many:media_tags" json:"tags" validate:"required"`
	Description   string    `json:"description" validate:"required"`
	AddedOn       string    `json:"addedOn"`
	Location      Location  `gorm:"embedded"`
	URL           string    `json:"url"`
	MimeType      EMimeType `json:"mimeType"`
}

type Tag struct {
	data.Identifiable
	Value string `json:"value" validate:"required"`
}

type SharedMedia struct {
	data.Identifiable
	Media []*Media `json:"media"`
}

type Story struct {
	data.Identifiable
	UserID        uint64      `json:"userId" gorm:"type:numeric"`
	SharedMedia   SharedMedia `json:"sharedMedia"`
	SharedMediaID uint64      `json:"sharedMediaId" gorm:"type:numeric"`
	CloseFriends  bool        `json:"closeFriends"`
}

type Post struct {
	data.Identifiable
	UserID        uint64      `json:"userId" gorm:"type:numeric"`
	SharedMedia   SharedMedia `validate:"required"`
	SharedMediaID uint64      `json:"sharedMediaId" gorm:"type:numeric"`
}

type ProfilePicture struct {
	data.Identifiable
	UserID uint64 `gorm:"type:numeric" json:"userId"`
	URL    string `json:"url"`
}

func PRToDataMedia(pr *prcontent.Media) *Media {
	tags := []Tag{}
	for _, t := range pr.Tags {
		tags = append(tags, Tag{Value: t.Value})
	}

	return &Media{
		SharedMediaID: pr.SharedMediaId,
		Filename:      pr.Filename,
		Tags:          tags,
		AddedOn:       pr.AddedOn,
		Description:   pr.Description,
		Location: Location{
			Country: pr.Location.Country,
			State:   pr.Location.State,
			ZipCode: pr.Location.ZipCode,
			Street:  pr.Location.Street,
		},
	}
}

func DataToPRStory(d *Story) *prcontent.Story {
	media := []*prcontent.Media{}
	for _, m := range d.SharedMedia.Media {
		media = append(media, DataToPRMedia(m))
	}

	return &prcontent.Story{
		Id:           d.ID,
		UserId:       d.UserID,
		CloseFriends: d.CloseFriends,
		Media:        media,
	}
}

func DataToPRMedia(d *Media) *prcontent.Media {
	tags := []*prcontent.Tag{}
	for _, t := range d.Tags {
		tags = append(tags, &prcontent.Tag{
			Value: t.Value,
			Id:    t.ID,
		})
	}

	mimeType := prcontent.EMimeType_IMAGE
	if d.MimeType == EMimeType_VIDEO {
		mimeType = prcontent.EMimeType_VIDEO
	}

	return &prcontent.Media{
		Id:          d.ID,
		Filename:    d.Filename,
		Description: d.Description,
		AddedOn:     d.AddedOn,
		Location: &prcontent.Location{
			Country: d.Location.Country,
			State:   d.Location.State,
			ZipCode: d.Location.ZipCode,
			Street:  d.Location.Street,
		},
		SharedMediaId: d.SharedMediaID,
		Tags:          tags,
		Url:           d.URL,
		MimeType:      mimeType,
	}
}

func (db *DBConn) GetSharedMediaByUser(id uint64) (*[]SharedMedia, error) {
	sharedMedia := []SharedMedia{}
	err := db.DB.Where("user_id = ?", id).Find(&sharedMedia).Error
	return &sharedMedia, err
}

// TODO(Jovan): REMOVE
func (m *Media) AfterCreate(tx *gorm.DB) error {
	logrus.Info("created")
	return nil
}

func (db *DBConn) AddMediaToSharedMedia(sharedMediaId uint64, media *Media) error {
	sm, err := db.GetSharedMedia(sharedMediaId)
	if err != nil {
		return err
	}
	return db.DB.Model(sm).Association("Media").Append(media)
}

func (db *DBConn) AddMediaToPost(postId uint64, media *Media) error {
	post, err := db.GetPost(postId)
	if err != nil {
		return err
	}
	return db.AddMediaToSharedMedia(post.SharedMediaID, media)
}

func (db *DBConn) AddMediaToStory(storyId uint64, media *Media) error {
	story, err := db.GetStory(storyId)
	if err != nil {
		return err
	}
	return db.AddMediaToSharedMedia(story.SharedMediaID, media)
}

var ErrMediaNotFound = fmt.Errorf("media not found")

func (db *DBConn) GetMediaByIds(ids ...uint64) ([]*Media, error) {
	media := []*Media{}
	res := db.DB.Preload(clause.Associations).Find(&media, ids)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, ErrMediaNotFound
	}

	return media, nil
}

var ErrStoryNotFound = fmt.Errorf("story not found")

func (db *DBConn) GetStory(storyId uint64) (*Story, error) {
	story := Story{}
	res := db.DB.Preload("SharedMedia").First(&story, storyId)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, ErrStoryNotFound
	}
	return &story, nil
}

func (db *DBConn) AddStory(s *Story) error {
	sharedMedia := SharedMedia{}
	if err := db.DB.Create(&sharedMedia).Error; err != nil {
		return err
	}
	s.SharedMediaID = sharedMedia.ID
	return db.DB.Create(s).Error
}

func (db *DBConn) GetSharedMedia(id uint64) (*SharedMedia, error) {
	sm := &SharedMedia{}
	res := db.DB.First(sm, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return sm, nil
}

var ErrStoriesNotFound = fmt.Errorf("stories not found")

func (db *DBConn) GetStoryByUser(id uint64) ([]*Story, error) {
	story := []*Story{}
	err := db.DB.Preload("SharedMedia.Media").Preload(clause.Associations).Where("user_id = ?", id).Find(&story).Error
	return story, err
}

func (db *DBConn) GetStoriesByUserAsMedia(userId uint64) ([]*Media, error) {
	media := []*Media{}
	stories, err := db.GetStoryByUser(userId)
	if err != nil {
		return nil, err
	}
	if len(stories) == 0 {
		return nil, ErrStoriesNotFound
	}
	for _, s := range stories {
		media = append(media, s.SharedMedia.Media...)
	}
	return media, nil
}

func (db *DBConn) GetPostByUser(id uint64) (*[]Post, error) {
	post := []Post{}
	err := db.DB.Preload("SharedMedia.Media").Preload(clause.Associations).Where("user_id = ?", id).Find(&post).Error
	return &post, err
}

var ErrPostNotFound = fmt.Errorf("post not found")

func (db *DBConn) GetPost(id uint64) (*Post, error) {
	post := Post{}
	res := db.DB.Preload("SharedMedia").First(&post, id)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, ErrPostNotFound
	}
	return &post, nil
}

func (db *DBConn) AddPost(p *Post) error {
	sharedMedia := SharedMedia{}
	if err := db.DB.Create(&sharedMedia).Error; err != nil {
		return err
	}
	p.SharedMediaID = sharedMedia.ID
	return db.DB.Create(p).Error
}

var ErrProfilePictureNotFound = fmt.Errorf("profile picture not found")

func (db *DBConn) GetProfilePictureByUser(id uint64) (*ProfilePicture, error) {
	post := ProfilePicture{}
	strId := strconv.FormatUint(id, 10)
	res := db.DB.Where("user_id = ?", strId).Find(&post)
	if res.RowsAffected == 0 {
		return nil, ErrProfilePictureNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &post, nil
}

func (db *DBConn) AddProfilePicture(pp *ProfilePicture) error {
	oldPP, err := db.GetProfilePictureByUser(pp.UserID)
	if err == ErrProfilePictureNotFound {
		return db.DB.Create(pp).Error
	}
	if err != nil {
		return err
	}
	oldPP.URL = pp.URL
	return db.DB.Save(oldPP).Error
}

func (db *DBConn) GetPostsByReaction(userId uint64) (*[]Post, error) {
	post := []Post{}
	err := db.DB.Preload("SharedMedia.Media").Preload(clause.Associations).Raw("SELECT p.* FROM posts p INNER JOIN reactions r on p.id = r.post_id WHERE r.user_id = ?", userId).Find(&post).Error
	return &post, err
}
