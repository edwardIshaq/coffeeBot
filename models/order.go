package models

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

// Order tracks the status of a Beverage order
type Order struct {
	gorm.Model
	BeverageID  uint `sql:"type:serial REFERENCES beverages(id)"`
	IsConfirmed bool
	IsFulfilled bool
}

// SaveNewOrder Creates an order in the DB
func SaveNewOrder(b Beverage) Order {
	order := Order{
		BeverageID:  b.ID,
		IsConfirmed: false,
		IsFulfilled: false,
	}
	db.Create(&order)
	return order
}

// OrderByID finds a beverage by ID (string)
func OrderByID(id string) Order {
	orderID, err := strconv.Atoi(id)
	if err != nil {
		return Order{}
	}
	order := Order{}
	db.First(&order, orderID)
	return order
}

// Confirm update the model to be confirmed
func (o Order) Confirm() {
	o.IsConfirmed = true
	db.Save(o)
}
