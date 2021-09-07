package data

type Campaign struct {
	Identifiable
	CampaignID uint64 `gorm:"type:numeric" json:"campaignId"`
}
