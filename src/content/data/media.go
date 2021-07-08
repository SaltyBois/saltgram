package data

import (
	"fmt"
	"saltgram/data"
	"saltgram/protos/content/prcontent"
	"strconv"

	"gorm.io/gorm/clause"
)

type EMimeType int32

const (
	EMimeType_IMAGE = iota
	EMimeType_VIDEO
)

type EAgeGroup int32

const (
	EAgeGroup_PRE20 = iota
	EAgeGroup_20
	EAgeGroup_30
)

type Media struct {
	data.Identifiable
	SharedMediaID uint64    `json:"sharedMediaId" gorm:"type:numeric"`
	Filename      string    `json:"filename" validate:"required"`
	Tags          []Tag     `gorm:"many2many:media_tags;" json:"tags" validate:"required"`
	Description   string    `json:"description" validate:"required"`
	AddedOn       string    `json:"addedOn"`
	Location      Location  `gorm:"embedded"`
	URL           string    `json:"url"`
	MimeType      EMimeType `json:"mimeType"`
	TaggedUsers   []UserTag `gorm:"many2many:media_taggedusers;" json:"taggedUsers"`
}

type Tag struct {
	data.Identifiable
	Value string `json:"value" validate:"required"`
}

type SharedMedia struct {
	data.Identifiable
	Media []*Media `json:"media"`
	// NOTE(Jovan): Flag for whether to read campaign
	// related fields
	IsCampaign       bool      `json:"isCampaign"`
	CampaignWebsite  string    `json:"campaignWebsite"`
	CampaignOneTime  bool      `json:"oneTime"`
	CampaignStart    string    `json:"campaignStart"`
	CampaignEnd      string    `json:"campaignEnd"`
	CampaignAgeGroup EAgeGroup `json:"ageGroup"`
	CampaignInfluencers []Influencer `gorm:"many2many:influencer_campaign;"`
}

type Influencer struct {
	data.Identifiable
	InfluencerID uint64 `gorm:"type:numeric"`

}

type CampaignChange struct {
	data.Identifiable
	CampaignID      int64  `json:"campaignID" gson:"type:numeric"`
	EffectiveAfter  string `json:"effectiveAfter"`
	CampaignOneTime bool   `json:"isCampaign"`
	CampaignWebsite string `json:"campaignWebsite"`
	CampaignEnd     string `json:"campaignEnd"`
}

type Story struct {
	data.Identifiable
	UserID        uint64      `json:"userId" gorm:"type:numeric"`
	SharedMedia   SharedMedia `json:"sharedMedia"`
	SharedMediaID uint64      `json:"sharedMediaId" gorm:"type:numeric"`
	CloseFriends  bool        `json:"closeFriends"`
}

type UserTag struct {
	UserID uint64 `gorm:"primaryKey; type:numeric" json:"userId"`
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

type SavedPost struct {
	UserID uint64 `json:"userId" gorm:"type:numeric"`
	PostID uint64 `json:"postId" gorm:"type:numeric"`
}

func PRToDataMedia(pr *prcontent.Media) *Media {
	tags := []Tag{}
	for _, t := range pr.Tags {
		tags = append(tags, Tag{Value: t.Value})
	}

	userTags := []UserTag{}
	for _, t := range pr.UserTags {
		userTags = append(userTags, UserTag{UserID: t.Id})
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
			City:    pr.Location.City,
			Street:  pr.Location.Street,
			Name:    pr.Location.Name,
		},
		TaggedUsers: userTags,
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
		IsCampaign: d.SharedMedia.IsCampaign,
		CampaignWebsite: d.SharedMedia.CampaignWebsite,
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

	userTags := []*prcontent.UserTag{}
	for _, t := range d.TaggedUsers {
		userTags = append(userTags, &prcontent.UserTag{
			Id: t.UserID,
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
			City:    d.Location.City,
			Street:  d.Location.Street,
			Name:    d.Location.Name,
		},
		SharedMediaId: d.SharedMediaID,
		Tags:          tags,
		Url:           d.URL,
		MimeType:      mimeType,
		UserTags:      userTags,
	}
}

func (db *DBConn) GetCampaignsByUser(id uint64) (*[]SharedMedia, error) {
	campaigns := []SharedMedia{}
	posts, err := db.GetPostByUser(id)
	if err != nil {
		return nil, err
	}
	stories, err := db.GetStoryByUser(id)
	if err != nil {
		return nil, err
	}
	for _, p := range *posts {
		if p.SharedMedia.IsCampaign {
			campaigns = append(campaigns, p.SharedMedia)
		}
	}
	for _, s := range stories {
		if s.SharedMedia.IsCampaign {
			campaigns = append(campaigns, s.SharedMedia)
		}
	}
	return &campaigns, err
}

func (db *DBConn) GetSharedMediaByUser(id uint64) (*[]SharedMedia, error) {
	sharedMedia := []SharedMedia{}
	err := db.DB.Where("user_id = ?", id).Find(&sharedMedia).Error
	return &sharedMedia, err
}

func (db *DBConn) AddInfluencerToCampaign(campaignId, influencerId uint64) error {
	sm, err := db.GetSharedMedia(campaignId)
	if err != nil {
		return err
	}
	in := Influencer{InfluencerID: influencerId}
	err = db.DB.Create(&in).Error
	if err != nil {
		return err
	}
	return db.DB.Model(sm).Association("CampaignInfluencers").Append(&in)
}

func (db *DBConn) AddMediaToSharedMedia(sharedMediaId uint64, media *Media) error {
	sm, err := db.GetSharedMedia(sharedMediaId)
	if err != nil {
		return err
	}
	media.SharedMediaID = sm.ID
	err = db.AddMedia(media)
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

func (db *DBConn) AddMedia(media *Media) error {
	return db.DB.Create(media).Error
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
	res := db.DB.Preload("Media").Preload("CampaignInfluencers").First(sm, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return sm, nil
}

var ErrStoriesNotFound = fmt.Errorf("stories not found")

func (db *DBConn) GetStoryByUser(id uint64) ([]*Story, error) {
	story := []*Story{}
	err := db.DB.
		Preload("SharedMedia",
			`(is_campaign AND campaign_one_time AND campaign_start = CAST(CURRENT_DATE as VARCHAR))
			OR (is_campaign AND not campaign_one_time AND CAST(campaign_start AS DATE) <= CURRENT_DATE
				AND CAST(campaign_end AS DATE) >= CURRENT_DATE)
			OR NOT is_campaign`).
		Preload("SharedMedia.Media.Tags").
		Preload("SharedMedia.Media.TaggedUsers").Preload(clause.Associations).
		Where("user_id = ?", id).Find(&story).Error
	if err != nil {
		return nil, err
	}
	for i, s := range story {
		if len(s.SharedMedia.Media) == 0 {
			story = append(story[:i], story[i + 1:]...)
		}
	}

	cs := []Story{}
	err = db.DB.Preload("SharedMedia").Preload("SharedMedia.Media.Tags").
		Preload("SharedMedia.Media.TaggedUsers").
		Model(&Story{}).Joins("INNER JOIN shared_media ON shared_media.id = stories.shared_media_id").
		Joins("INNER JOIN influencer_campaign ic ON shared_media.id = ic.shared_media_id").
		Joins("INNER JOIN influencers i ON i.id = ic.influencer_id").
		Where("i.influencer_id = ?", id).Find(&cs).Error
	if err != nil {
		return nil, err
	}
	for _, s := range cs {
		if len(s.SharedMedia.Media) != 0 {
			story = append(story, &s)
		}
	}
	return story, nil
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
	err := db.DB.
	Preload("SharedMedia",
		`(is_campaign AND campaign_one_time AND campaign_start = CAST(CURRENT_DATE as VARCHAR))
		OR (is_campaign AND not campaign_one_time AND CAST(campaign_start AS DATE) <= CURRENT_DATE
			AND CAST(campaign_end AS DATE) >= CURRENT_DATE)
		OR NOT is_campaign`).
	Preload("SharedMedia.Media.Tags").Preload("SharedMedia.Media.TaggedUsers").Preload(clause.Associations).Where("user_id = ?", id).Find(&post).Error
	if err != nil {
		return nil, err
	}
	for i, p := range post {
		if len(p.SharedMedia.Media) == 0 {
			post = append(post[:i], post[i + 1:]...)
		}
	}
	cp := []Post{}
	err = db.DB.Preload("SharedMedia").Preload("SharedMedia.Media.Tags").
	Preload("SharedMedia.Media.TaggedUsers").
	Model(&Post{}).Joins("INNER JOIN shared_media ON shared_media.id = posts.shared_media_id").
	Joins("INNER JOIN influencer_campaign ic ON shared_media.id = ic.shared_media_id").
	Joins("INNER JOIN influencers i ON i.id = ic.influencer_id").
	Where("i.influencer_id = ?", id).Find(&cp).Error
	if err != nil {
		return nil, err
	}
	for _, p := range cp {
		if len(p.SharedMedia.Media) != 0 {
			post = append(post, p)
		}
	}
	return &post, err
}

var ErrPostNotFound = fmt.Errorf("post not found")

func (db *DBConn) GetPost(id uint64) (*Post, error) {
	post := Post{}
	res := db.DB.Preload("SharedMedia.Media.Tags").Preload("SharedMedia.Media.TaggedUsers").Preload(clause.Associations).First(&post, id)
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
	err := db.DB.Preload("SharedMedia.Media.Tags").Preload("SharedMedia.Media.TaggedUsers").Preload(clause.Associations).Raw("SELECT p.* FROM posts p INNER JOIN reactions r on p.id = r.post_id WHERE r.user_id = ?", userId).Find(&post).Error
	return &post, err
}

func (db *DBConn) AddTag(t *Tag) (*Tag, error) {
	tag := t
	err := db.DB.Create(tag).Error
	fmt.Println(tag.ID)
	return tag, err
}

func (db *DBConn) GetTagById(id uint64) (Tag, error) {
	tag := Tag{}
	err := db.DB.Where("id = ?", id).First(&tag).Error
	return tag, err
}

var ErrTagNotExists = fmt.Errorf("tag not exists")

func (db *DBConn) GetIfExists(value string) (*Tag, error) {
	tag := Tag{}
	res := db.DB.Where("value = ?", value).First(&tag)
	if res.RowsAffected == 0 {
		return nil, ErrTagNotExists
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &tag, res.Error
}

func (db *DBConn) AddUserTag(t *UserTag) (*UserTag, error) {
	tag := t
	err := db.DB.Create(tag).Error
	return tag, err
}

var ErrUserTagNotFound = fmt.Errorf("user tag not found")

func (db *DBConn) GetUserTagById(id uint64) (*UserTag, error) {
	userTag := UserTag{}
	res := db.DB.Where("user_id = ?", id).First(&userTag)
	if res.RowsAffected == 0 {
		return nil, ErrUserTagNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &userTag, res.Error
}

var ErrTagNotFound = fmt.Errorf("tag not found")

func (db *DBConn) GetTagId(value string) (uint64, error) {
	tag := Tag{}
	res := db.DB.Where("value = ?", value).First(&tag)
	if res.RowsAffected == 0 {
		return 0, ErrTagNotFound
	}
	if res.Error != nil {
		return 0, res.Error
	}
	return tag.ID, res.Error
}

var ErrSharedMediaNotFound = fmt.Errorf("shared media id not found")

func (db *DBConn) GetSharedMediaIdByTagId(tagId uint64) ([]uint64, error) {
	var id []uint64
	res := db.DB.Raw("select distinct m.shared_media_id from media m inner join media_tags mt on m.id = mt.media_id where mt.tag_id = ?", tagId).Find(&id)
	if res.RowsAffected == 0 {
		return nil, ErrSharedMediaNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return id, res.Error
}

func (db *DBConn) GetPostsBySharedMediaId(ids []uint64) (*[]Post, error) {
	posts := []Post{}
	res := db.DB.Preload("SharedMedia.Media.Tags").Preload("SharedMedia.Media.TaggedUsers").Preload(clause.Associations).Where("shared_media_id IN ?", ids).Find(&posts)
	return &posts, res.Error
}

func (db *DBConn) GetPostsByTag(value string) (*[]Post, error) {
	id, err := db.GetTagId(value)
	if err != nil {
		return nil, err
	}

	ids, err := db.GetSharedMediaIdByTagId(id)
	if err != nil {
		return nil, err
	}

	posts, err := db.GetPostsBySharedMediaId(ids)
	return posts, err
}

func (db *DBConn) GetAllTagsByNameSubstring(value string) ([]Tag, error) {
	var tags []Tag
	query := "%" + value + "%"
	err := db.DB.Where("value LIKE ?", query).Limit(21).Find(&tags).Error
	return tags, err
}

func (db *DBConn) GetAllLocationNames(name string) ([]string, error) {
	var names []string
	query := "%" + name + "%"
	err := db.DB.Raw("SELECT DISTINCT name FROM media WHERE name LIKE ?", query).Limit(21).Find(&names).Error
	return names, err
}

func (db *DBConn) GetSharedMediaIdByLocationName(name string) ([]uint64, error) {
	var id []uint64
	res := db.DB.Raw("select distinct shared_media_id from media where name = ?", name).Find(&id)
	if res.RowsAffected == 0 {
		return nil, ErrSharedMediaNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return id, res.Error
}

func (db *DBConn) GetContentsByLocation(name string) (*[]Post, error) {
	ids, err := db.GetSharedMediaIdByLocationName(name)
	if err != nil {
		return nil, err
	}

	posts, err := db.GetPostsBySharedMediaId(ids)
	return posts, err
}

func (db *DBConn) AddSavedPost(sp *SavedPost) error {
	savedPost := sp
	return db.DB.Create(savedPost).Error
}

func (db *DBConn) GetSavedPosts(userId uint64) (*[]Post, error) {
	var ids []uint64
	res := db.DB.Raw("select distinct post_id from saved_posts where user_id = ?", userId).Find(&ids)
	if res.Error != nil {
		return nil, res.Error
	}
	posts := []Post{}
	err := db.DB.Preload("SharedMedia.Media.Tags").Preload("SharedMedia.Media.TaggedUsers").Preload(clause.Associations).Where("id IN ?", ids).Find(&posts).Error
	return &posts, err
}

func (db *DBConn) GetTaggedPostsByUser(userId uint64) (*[]Post, error) {
	var ids []uint64
	res := db.DB.Raw("select distinct m.shared_media_id from media m inner join media_taggedusers mt on m.id = mt.media_id where mt.user_tag_user_id = ?", userId).Find(&ids)
	if res.Error != nil {
		return nil, res.Error
	}

	posts, err := db.GetPostsBySharedMediaId(ids)
	return posts, err
}
