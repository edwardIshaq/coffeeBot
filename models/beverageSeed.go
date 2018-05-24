package models

/*
This will provide the base beverages to seed the DB
*/

func makeCoffee(name, espresso, cup, milk string) Beverage {
	return Beverage{
		DefaultDrink:   true,
		Category:       "Coffee",
		Name:           name,
		Espresso:       espresso,
		CupType:        cup,
		Syrup:          syrupNone,
		Temperture:     tempHot,
		DrinkOfTheWeek: false,
	}
}

// CoffeeDrinks basic coffee stuff
func CoffeeDrinks() []Beverage {
	return []Beverage{
		makeCoffee("Espresso", espressoDouble, cupSize6oz, milkNone),
		makeCoffee("Macchiato", espressoSingle, cupSize8oz, milkNone),
		makeCoffee("Gibraltar / Cortado", espressoQuad, cupSize6oz, milkNone),
		makeCoffee("Cappuccino", espressoDouble, cupSize6oz, milk2),
	}
}

func makeTea(name, milk, syrup string) Beverage {
	return Beverage{
		DefaultDrink: true,
		Category:     "Tea",
		Name:         name,
		CupType:      cupSize8oz,
		Syrup:        syrup,
		Temperture:   tempHot,
	}
}

// TeaDrinks returns all standard `Tea` drinks
func TeaDrinks() []Beverage {
	return []Beverage{
		makeTea("London Fog", milkAlmond, syrupMint),
		makeTea("San-Fran Fog", milkSoy, syrupHoney),
		makeTea("Matcha Latte", milkAlmond, syrupMint),
		makeTea("Tanglewood Ginger Chai", milkCashew, syrupLavender),
	}
}

func makeDrinkOfTheWeek(name, espresso, cup, milk string) Beverage {
	return Beverage{
		DefaultDrink:   true,
		Category:       "Drink of the week",
		Name:           name,
		Espresso:       espresso,
		CupType:        cup,
		Syrup:          syrupNone,
		Temperture:     tempHot,
		DrinkOfTheWeek: true,
	}
}

// DrinkOfTheWeek list of the drink of the week
func DrinkOfTheWeek() []Beverage {
	return []Beverage{
		makeDrinkOfTheWeek("Vitality Latte", espressoNone, cupSize6oz, milkNone),
		makeDrinkOfTheWeek("Herbal Remedy Tea", espressoNone, cupSize6oz, milkNone),
		makeDrinkOfTheWeek("Iced Separator", espressoNone, cupSize6oz, milkNone),
	}
}
