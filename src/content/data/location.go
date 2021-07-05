package data

type Location struct {
	Country string `json:"country"`
	State   string `json:"state"`
	ZipCode string `json:"zipcode"`
	City    string `json:"city"`
	Street  string `json:"street"`
	Name    string `json:"name"`
}
