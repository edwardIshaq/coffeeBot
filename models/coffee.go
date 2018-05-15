package models

import (
	"github.com/jinzhu/gorm"
)

// cupType is the drink size of a beverage
type cupType string

const (
	cupTypeGibraltar cupType = "Gibraltar Glass"
	cupSize6oz       cupType = "6oz cup"
	cupSize8oz       cupType = "8oz cup"
	cupSize12oz      cupType = "12oz mug"
	cupSize8ozPaper  cupType = "8oz paper"
	cupSize12ozPaper cupType = "12oz paper"
	cupSizePint      cupType = "Pint glass"
	cupSizeDefault   cupType = "regular"
)

// AllDrinkSizes a list of all possible DrinkSize
func AllDrinkSizes() []string {
	return []string{
		string(cupTypeGibraltar),
		string(cupSize6oz),
		string(cupSize8oz),
		string(cupSize12oz),
		string(cupSize8ozPaper),
		string(cupSize12ozPaper),
		string(cupSizePint),
		string(cupSizeDefault),
	}
}

// Temperture of the beverage
type Temperture string

const (
	tempHot  Temperture = "hot"
	tempIced Temperture = "iced"
)

// AllTemps all temps
func AllTemps() []string {
	return []string{
		string(tempHot),
		string(tempIced),
	}
}

// espressoOption to captuer the espresso shots in a beverage
type espressoOption string

const (
	espressoNone    espressoOption = "None"
	espressoSingle  espressoOption = "single"
	espressoDouble  espressoOption = "double"
	espressoTripple espressoOption = "triple"
	espressoQuad    espressoOption = "quad"
	espressoDecaf   espressoOption = "decaf"
	espressoHalfCaf espressoOption = "Half-caf"
)

// AllEspressoOptions a list of all possible DrinkSize
func AllEspressoOptions() []string {
	return []string{
		string(espressoNone),
		string(espressoSingle),
		string(espressoDouble),
		string(espressoTripple),
		string(espressoQuad),
		string(espressoDecaf),
		string(espressoHalfCaf),
	}
}

// syrup place holder type
type syrup string

const (
	syrupNone      syrup = "No Syrup"
	syrupMint      syrup = "Mint"
	syrupLavender  syrup = "Lavender"
	syrupChocolate syrup = "Chocolate"
	syrupVanilla   syrup = "Vanilla"
	syrupHoney     syrup = "Honey"
)

// AllSyrupOptions a list of all the Syrup options
func AllSyrupOptions() []string {
	return []string{
		string(syrupNone),
		string(syrupMint),
		string(syrupLavender),
		string(syrupChocolate),
		string(syrupVanilla),
		string(syrupHoney),
	}
}

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

func newBev() Beverage {
	bev := Beverage{}
	bev.CupType = string(cupSizeDefault)
	bev.Temperture = string(tempHot)
	return bev
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
	return defaultBeverage(name)
}

func defaultBeverage(name string) Beverage {
	return Beverage{
		Espresso:   string(espressoSingle),
		Syrup:      string(syrupMint),
		Temperture: string(tempHot),
		CupType:    string(cupSize8oz),
	}
}

// Drink Menus

// AllDrinksOfTheWeek list
func AllDrinksOfTheWeek() []string {
	return []string{"Vitality Latte", "Herbal Remedy Tea", "Iced Separator"}
}

// AllUsualDrinks list
func AllUsualDrinks() []string {
	return []string{"Steamed Milk", "Hot Chocolate", "Tea"}
}

// AllTeas list
func AllTeas() []string {
	return []string{"London Fog", "San-Fran Fog", "Matcha Latte", "Tanglewood Ginger Chai"}
}

// AllCoffees list
func AllCoffees() []string {
	return []string{"Espresso", "Macchiato", "Gibraltar / Cortado", "Cappuccino"}
}
