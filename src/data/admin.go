package data

type InappropriateContentReportDTO struct {
	SharedMediaId uint64 `json:"sharedMediaId"`
	UserId        uint64 `json:"userId"`
}

type ReviewRequestDTO struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type VerificationRequestDTO struct {
	Id       string `json:"id"`
	FullName string `json:"fullname"`
	Category string `json:"category"`
	Url      string `json:"url"`
	UserId   uint64 `json:"userId"`
	Username string `json:"username"`
}
