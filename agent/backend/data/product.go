package data

import "fmt"

type Product struct {
	Identifiable
	ImageURL string  `json:"imageURL"`
	Price    float32 `json:"price"`
	Amount   int32   `json:"amount"`
}

func (db *DBConn) AddProduct(p *Product) error {
	return db.DB.Create(p).Error
}

var ErrProductNotFound = fmt.Errorf("product not found")

func (db *DBConn) UpdateProduct(p *Product) error {
	prod := Product{}
	res := db.DB.First(&prod, p.ID)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrProductNotFound
	}
	prod.Amount = p.Amount
	prod.ImageURL = p.ImageURL
	prod.Price = p.Price
	return db.DB.Save(&prod).Error
}

func (db *DBConn) RemoveProduct(id uint64) error {
	prod := Product{}
	prod.ID = id
	return db.DB.Delete(&prod).Error
}

func (db *DBConn) GetAll() (*[]Product, error) {
	prods := []Product{}
	err := db.DB.Find(&prods).Error
	return &prods, err
}
