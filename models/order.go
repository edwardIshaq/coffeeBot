package models

import "github.com/jinzhu/gorm"

// Order tracks the status of a Beverage order
type Order struct {
	gorm.Model
	BeverageID  uint `sql:"type:serial REFERENCES beverages(id)"`
	IsConfirmed bool
	IsFulfilled bool
}

// SaveNewOrder Creates an order in the DB
func SaveNewOrder(b Beverage) {
	order := Order{
		BeverageID:  b.ID,
		IsConfirmed: false,
		IsFulfilled: false,
	}
	db.Create(&order)
}
