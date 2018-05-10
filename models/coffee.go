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

func newBeverage(name string, espresso espressoOption, syrup syrup) Beverage {
	return Beverage{
		Name:       name,
		Espresso:   string(espresso),
		Syrup:      string(syrup),
		Temperture: string(tempHot),
		CupType:    string(cupSize12oz),
	}
}

func beverageList() []Beverage {
	espresso := newBeverage("Espresso", espressoSingle, syrupChocolate)
	espresso.CupType = string(cupSize8oz)

	return []Beverage{
		newBeverage("Vitality Latte", espressoQuad, syrupHoney),
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
	// return Beverage{Name: name}
	return demoBevMatch(name)
}

func demoBevMatch(name string) Beverage {
	return beverageList()[0]
}
