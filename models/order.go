package models

import "github.com/jinzhu/gorm"

// Order tracks the status of a Beverage order
type Order struct {
	gorm.Model
	BeverageID  int `sql:"type:serial REFERENCES beverages(id)"`
	IsConfirmed bool
	IsFulfilled bool
}
