package models

/*
https://slackcoffeebar.typeform.com/to/v6kODV

TODO:
[ ] Add all items to DB

*/

// cupType is the drink size of a beverage
const (
	cupSize6oz       string = "6oz cup"
	cupSize8oz       string = "8oz cup"
	cupSize12oz      string = "12oz mug"
	cupSize8ozPaper  string = "8oz paper"
	cupSize12ozPaper string = "12oz paper"
	cupSizePint      string = "Pint glass"
	cupSizeGibraltar string = "Gibraltar Glass"
	cupSizeDefault   string = "regular"
)

// AllDrinkSizes a list of all possible DrinkSize
func AllDrinkSizes() []string {
	return []string{
		cupSize6oz,
		cupSize8oz,
		cupSize12oz,
		cupSize8ozPaper,
		cupSize12ozPaper,
		cupSizePint,
		cupSizeGibraltar,
		cupSizeDefault,
	}
}

// Temperture of the beverage
const (
	tempHot  string = "hot"
	tempIced string = "iced"
)

// AllTemps all temps
func AllTemps() []string {
	return []string{
		tempHot,
		tempIced,
	}
}

// espressoOption to captuer the espresso shots in a beverage
const (
	espressoNone    string = "None"
	espressoSingle  string = "single"
	espressoDouble  string = "double"
	espressoTripple string = "triple"
	espressoQuad    string = "quad"
	espressoDecaf   string = "decaf"
	espressoHalfCaf string = "Half-caf"
)

// AllEspressoOptions a list of all possible DrinkSize
func AllEspressoOptions() []string {
	return []string{
		espressoNone,
		espressoSingle,
		espressoDouble,
		espressoTripple,
		espressoQuad,
		espressoDecaf,
		espressoHalfCaf,
	}
}

// syrup place holder type
const (
	syrupNone      string = "No Syrup"
	syrupMint      string = "Mint"
	syrupLavender  string = "Lavender"
	syrupChocolate string = "Chocolate"
	syrupVanilla   string = "Vanilla"
	syrupHoney     string = "Honey"
)

// AllSyrupOptions a list of all the Syrup options
func AllSyrupOptions() []string {
	return []string{
		syrupNone,
		syrupMint,
		syrupLavender,
		syrupChocolate,
		syrupVanilla,
		syrupHoney,
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
