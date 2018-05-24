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
	saveAllDrinksToDB()
}

func setupTables() {
	db.AutoMigrate(&Team{})
	db.AutoMigrate(&Beverage{})
}

// a hack to prevent duplicating the default drinks
func saveAllDrinksToDB() {
	drinks := CoffeeDrinks()
	drinks = append(drinks, TeaDrinks()...)
	drinks = append(drinks, DrinkOfTheWeek()...)

	var count int
	db.Where(&Beverage{DefaultDrink: true}).Model(&Beverage{}).Count(&count)
	if count == len(drinks) {
		return
	}

	for _, bev := range drinks {
		db.Save(&bev)
	}
}
