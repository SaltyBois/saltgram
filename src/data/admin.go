package data

/*type Category string

const (
	INFLUENCER Category = "INFLUENCER"
	SPORTS
	MEDIA
	BUSINESS
	BRAND
	ORGANIZATION
)
*/
type VerificationRequestDTO struct {
	FullName      string `json:"fullname" validate:"required"`
	Category      string `json:"category" validate:"required"`
	UserId        uint64
	DocumentImage MediaDTO
}

type InappropriateContentReportDTO struct {
	SharedMedia SharedMediaDTO `json:"sharedMedia"`
	User        UserDTO        `json:"user"`
	Reason      string         `json:"reason"`
}
