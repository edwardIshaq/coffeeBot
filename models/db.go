package models

import (
	"github.com/jinzhu/gorm"
)

// Private package varible, injected in by calling `SetDatabase`
var db *gorm.DB

// SetDatabase sets up the DB
func SetDatabase(database *gorm.DB) {
	db = database
	db.LogMode(true)
	setupTables()
}

func setupTables() {
	db.AutoMigrate(&Team{})
	db.AutoMigrate(&Beverage{})
}
