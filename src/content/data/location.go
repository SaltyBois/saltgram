package data

type Location struct {
	Id       string     `json:"id" gorm:"primaryKey" validate:"required"`
	Country  string     `json:"country" validate:"required"`
	State    string      `json:"state" validate:"required"`
	ZipCode  string       `json:"zipcode" validate:"required"`
	Street   string       `json:"street" validate:"required"`
}


