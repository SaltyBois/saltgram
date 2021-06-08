package data

type Location struct {
	Country string `json:"country" validate:"required"`
	State   string `json:"state" validate:"required"`
	ZipCode string `json:"zipcode" validate:"required"`
	Street  string `json:"street" validate:"required"`
}
