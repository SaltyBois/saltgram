package data

import (
	"saltgram/protos/content/prcontent"

	"github.com/go-playground/validator"
)

type SharedMediaDTO struct {
	Media []MediaDTO
}

type MediaDTO struct {
	Id          string   `json:"id"`
	UserId      uint64
	Filename    string   `json:"filename" validate:"required"`
	Tags        []TagDTO `gorm:"many2many:media_tags" json:"tags" validate:"required"`
	Description string   `json:"description" validate:"required"`
	AddedOn     string   `json:"addedOn"`
	Location    LocationDTO
	SharedMediaID string `json:"sharedMediaId"`
	URL         string `json:"url"`
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
	SharedMedia SharedMediaDTO
	User        UserDTO
}

type CommentDTO struct {
	Content string
	UserId  uint64
	PostId  uint64
}

type ReactionType string

const (
	LIKE ReactionType = "LIKE"
	DISLIKE
)

type ReactionDTO struct {
	ReactionType ReactionType
	UserId       uint64
	PostId       uint64
}

type HighlightRequest struct {
	Name string `json:"name"`
	Stories []MediaDTO `json:"stories"`
}

func (sm *SharedMediaDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(sm)
}

func PRToDTOTag(pr *prcontent.Tag) *TagDTO {
	return &TagDTO{
		ID: pr.Id,
		Value: pr.Value,
	}
}

func PRToDTOLocation(pr *prcontent.Location) *LocationDTO {
	return &LocationDTO{
		Country: pr.Country,
		State: pr.State,
		ZipCode: pr.ZipCode,
		Street: pr.Street,
	}
}

func DTOToPRLocation(dto *LocationDTO) *prcontent.Location {
	return &prcontent.Location{
		Country: dto.Country,
		State: dto.State,
		ZipCode: dto.ZipCode,
		Street: dto.Street,
	}
}

func DTOToPRTag(dto *TagDTO) *prcontent.Tag {
	return &prcontent.Tag{
		Value: dto.Value,
		Id: dto.ID,
	}
}