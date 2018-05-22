package models

/*
https://slackcoffeebar.typeform.com/to/v6kODV

TODO:
[ ] Add all items to DB

*/

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
