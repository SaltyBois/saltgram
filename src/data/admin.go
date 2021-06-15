package data

type Category string

const (
	INFLUENCER Category = "INFLUENCER"
	SPORTS
	MEDIA
	BUSINESS
	BRAND
	ORGANIZATION
)

type VerificationRequestDTO struct {
	FullName string   `json:"fullname" validate:"required"`
	Category Category `json:"category" validate:"required"`
	UserId   uint64
	//DocumentImage
}
