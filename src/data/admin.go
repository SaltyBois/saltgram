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

type InappropriateContentReportDTO struct {
	SharedMedia SharedMediaDTO `json:"sharedMedia"`
	User        UserDTO        `json:"user"`
}

type ReviewRequestDTO struct {
	Id     uint64 `json:"id"`
	Status string `json:"status"`
}

type VerificationRequestDTO struct {
	FullName string `json:"fullname"`
	Category string `json:"category"`
	Url      string `json:"url"`
}
