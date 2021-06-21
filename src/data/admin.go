package data

type InappropriateContentReportDTO struct {
	SharedMediaId string `json:"sharedMediaId"`
	UserId        uint64 `json:"userId"`
}

type GetInappropriateContentReportDTO struct {
	Id             string `json:"id"`
	SharedMediaId  string `json:"sharedMediaId"`
	UserId         uint64 `json:"userId"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profilePicture"`
}

type ReviewRequestDTO struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type VerificationRequestDTO struct {
	Id             string `json:"id"`
	FullName       string `json:"fullname"`
	Category       string `json:"category"`
	Url            string `json:"url"`
	UserId         uint64 `json:"userId"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profilePicture"`
}
