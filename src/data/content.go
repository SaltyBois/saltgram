package data

import "github.com/go-playground/validator"

type SharedMediaDTO struct {
	Media []MediaDTO
}

type MediaDTO struct {
	UserId      uint64
	Filename    string   `json:"filename" validate:"required"`
	Tags        []TagDTO `gorm:"many2many:media_tags" json:"tags" validate:"required"`
	Description string   `json:"description" validate:"required"`
	AddedOn     string   `json:"addedOn"`
	Location    LocationDTO
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

func (sm *SharedMediaDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(sm)
}