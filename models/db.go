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
	db.Model(&Beverage{}).AddUniqueIndex("uniq_idx_default_bevs", "category", "name", "default_drink")
}

// a hack to prevent duplicating the default drinks
func saveAllDrinksToDB() {
	drinks := CoffeeDrinks()
	drinks = append(drinks, TeaDrinks()...)
	drinks = append(drinks, DrinkOfTheWeek()...)

	for _, bev := range drinks {
		db.Save(&bev)
	}
}
