package models

import "github.com/jinzhu/gorm"

// Beverage type will hold the default selections for the menus
type Beverage struct {
	gorm.Model
	Name            string
	Espresso        string
	Syrup           string
	Temperture      string
	CupType         string
	DrinkOfTheWeek  bool
	BaristaApproved bool
	Comment         string
}

func newBeverage(name string, espresso espressoOption, syrup syrup, cup cupType) Beverage {
	return Beverage{
		Name:       name,
		Espresso:   string(espresso),
		Syrup:      string(syrup),
		Temperture: string(tempHot),
		CupType:    string(cup),
	}
}

func beverageList() []Beverage {
	espresso := newBeverage("Espresso", espressoSingle, syrupNone, cupSize8oz)

	return []Beverage{
		espresso,
		newBeverage("Vitality Latte", espressoQuad, syrupHoney, cupSize8oz),
	}
}

// BeverageByName gets you a preset beverage or an empty one if not found
func BeverageByName(name string) Beverage {
	list := beverageList()
	for _, bev := range list {
		if bev.Name == name {
			return bev
		}
	}
	return newBeverage("", espressoSingle, syrupMint, cupSize8oz)
}
