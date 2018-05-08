package models

// DrinkSize is the drink size of a beverage
type DrinkSize string

const (
	drinkSize8oz     DrinkSize = "8 Oz"
	drinkSize12oz    DrinkSize = "12 Oz"
	drinkSizeDefault DrinkSize = "regular"
)

// AllDrinkSizes a list of all possible DrinkSize
func AllDrinkSizes() []string {
	return []string{"8 Oz", "12 Oz", "regular"}
}

// Temperture of the beverage
type Temperture string

const (
	tempHot  Temperture = "hot"
	tempIced Temperture = "iced"
)

// EspressoOption to captuer the espresso shots in a beverage
type EspressoOption string

const (
	// EspressoSingle single shot
	EspressoSingle EspressoOption = "single"
	// EspressoDouble double shot
	EspressoDouble EspressoOption = "double"
	// EspressoTripple triple shot
	EspressoTripple EspressoOption = "triple"
	// EspressoQuad quat shot
	EspressoQuad EspressoOption = "quad"
	// EspressoDecaf decaf
	EspressoDecaf EspressoOption = "decaf"
	// EspressoHalfCaf half-caf
	EspressoHalfCaf EspressoOption = "Half-caf"
)

// AllEspressoOptions a list of all possible DrinkSize
func AllEspressoOptions() []string {
	return []string{"single", "double", "triple", "quad", "decaf", "Half-caf"}
}

// Syrup place holder type
type Syrup string

const (
	//Mint Syrup
	Mint Syrup = "Mint"
	//Lavender Syrup
	Lavender Syrup = "Lavender"
	//Chocolate Syrup
	Chocolate Syrup = "Chocolate"
	//Vanilla Syrup
	Vanilla Syrup = "Vanilla"
	//Honey Syrup
	Honey Syrup = "Honey"
)

// AllSyrupOptions a list of all the Syrup options
func AllSyrupOptions() []string {
	return []string{"Mint", "Lavender", "Chocolate", "Vanilla", "Honey"}
}
