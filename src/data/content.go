package data

import (
	"saltgram/protos/content/prcontent"
	"strconv"

	"github.com/go-playground/validator"
)

type SharedMediaDTO struct {
	Media []MediaDTO
}

type MediaDTO struct {
	Id            string `json:"id"`
	UserId        uint64
	Filename      string   `json:"filename" validate:"required"`
	Tags          []TagDTO `gorm:"many2many:media_tags" json:"tags" validate:"required"`
	Description   string   `json:"description" validate:"required"`
	AddedOn       string   `json:"addedOn"`
	Location      LocationDTO
	SharedMediaID string `json:"sharedMediaId"`
	URL           string `json:"url"`
	MimeType      string `json:"mimeType"`
}

type LocationDTO struct {
	Country string `json:"country" validate:"required"`
	State   string `json:"state" validate:"required"`
	ZipCode string `json:"zipcode" validate:"required"`
	Street  string `json:"street" validate:"required"`
}

type TagDTO struct {
	ID    uint64 `json:"id"`
	Value string `json:"value" validate:"required"`
}

type PostDTO struct {
	Id          string `json:"id"`
	SharedMedia SharedMediaDTO
	UserId      string `json:"userId"`
}

type CommentDTO struct {
	Content string
	UserId  uint64
	PostId  string
}

type ReactionType string

const (
	LIKE ReactionType = "LIKE"
	DISLIKE
)

type ReactionDTO struct {
	ReactionType string
	UserId       uint64
	PostId       string
}

type HighlightRequest struct {
	Name    string     `json:"name"`
	Stories []MediaDTO `json:"stories"`
}

type HighlightDTO struct {
	Name    string     `json:"name"`
	Stories []MediaDTO `json:"stories"`
}

func (sm *SharedMediaDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(sm)
}

type ReactionPutDTO struct {
	Id           string
	ReactionType string
}

func PRToDTOHighlight(pr *prcontent.Highlight) *HighlightDTO {
	media := []MediaDTO{}
	for _, s := range pr.Stories {
		media = append(media, *PRToDTOMedia(s))
	}

	return &HighlightDTO{
		Name:    pr.Name,
		Stories: media,
	}
}

func PRToDTOMedia(pr *prcontent.Media) *MediaDTO {
	tags := []TagDTO{}
	for _, t := range pr.Tags {
		tags = append(tags, *PRToDTOTag(t))
	}

	id := strconv.FormatUint(pr.Id, 10)
	sharedMediaId := strconv.FormatUint(pr.SharedMediaId, 10)

	mimeType := "image"
	if pr.MimeType == prcontent.EMimeType_VIDEO {
		mimeType = "video"
	}

	return &MediaDTO{
		Id:            id,
		UserId:        pr.UserId,
		Filename:      pr.Filename,
		Description:   pr.Description,
		Tags:          tags,
		AddedOn:       pr.AddedOn,
		Location:      *PRToDTOLocation(pr.Location),
		SharedMediaID: sharedMediaId,
		URL:           pr.Url,
		MimeType:      mimeType,
	}
}

func PRToDTOTag(pr *prcontent.Tag) *TagDTO {
	return &TagDTO{
		ID:    pr.Id,
		Value: pr.Value,
	}
}

func PRToDTOLocation(pr *prcontent.Location) *LocationDTO {
	return &LocationDTO{
		Country: pr.Country,
		State:   pr.State,
		ZipCode: pr.ZipCode,
		Street:  pr.Street,
	}
}

func DTOToPRLocation(dto *LocationDTO) *prcontent.Location {
	return &prcontent.Location{
		Country: dto.Country,
		State:   dto.State,
		ZipCode: dto.ZipCode,
		Street:  dto.Street,
	}
}

func DTOToPRTag(dto *TagDTO) *prcontent.Tag {
	return &prcontent.Tag{
		Value: dto.Value,
		Id:    dto.ID,
	}
}
