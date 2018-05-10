package models

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
	cups := []cupType{cupTypeGibraltar, cupSize6oz,
		cupSize8oz, cupSize12oz, cupSize8ozPaper,
		cupSize12ozPaper, cupSizePint, cupSizeDefault}

	result := make([]string, len(cups))
	for idx, cup := range cups {
		result[idx] = string(cup)
	}
	return result
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
	// EspressoNone for none
	EspressoNone EspressoOption = "None"
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
	// SyrupNone for no Syrup
	SyrupNone Syrup = ""
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

// Beverage type will hold the default selections for the menus
type Beverage struct {
	Name       string
	Espresso   string
	Syrup      string
	Temperture string
	CupType    string
	Comment    string
}

func newBev() Beverage {
	bev := Beverage{}
	bev.CupType = string(cupSizeDefault)
	bev.Temperture = string(tempHot)
	return bev
}

func newBeverage(name string, espresso EspressoOption, syrup Syrup) Beverage {
	return Beverage{
		Name:       name,
		Espresso:   string(espresso),
		Syrup:      string(syrup),
		Temperture: string(tempHot),
	}
}

func beverageList() []Beverage {
	espresso := newBeverage("Espresso", EspressoSingle, SyrupNone)
	espresso.CupType = string(cupSize8oz)

	return []Beverage{
		newBeverage("Vitality Latte", EspressoNone, Honey),
		espresso,
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
	return Beverage{Name: name}
}
