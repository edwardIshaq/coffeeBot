package models

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
)

// Beverage type will hold the default selections for the menus
type Beverage struct {
	gorm.Model
	Name           string
	Espresso       string
	Syrup          string
	Temperture     string
	CupType        string
	DrinkOfTheWeek bool
	DefaultDrink   bool
	Comment        string
	UserID         string
}

func newBeverage(name, espresso, syrup, cup string) Beverage {
	return Beverage{
		Name:       name,
		Espresso:   string(espresso),
		Syrup:      string(syrup),
		Temperture: string(tempHot),
		CupType:    string(cup),
	}
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

/*
map[Comment: CupType:8oz paper Espresso:single Syrup:Lavender Temperture:hot]
*/
func saveBeverage(submission map[string]string, userID string) {
	bev := &Beverage{}
	var err error
	for k, v := range submission {
		err = setField(bev, k, v)
		if err != nil {
			break
		}
	}

	if err != nil {
		fmt.Printf("\nError creating `Beverage`: %v", err)
		return
	}

	bev.UserID = userID
	bev.DefaultDrink = false
	db.Save(&bev)
}

// BeverageByName gets you a preset beverage or an empty one if not found
func BeverageByName(name string) Beverage {
	for _, bev := range DefaultDrinks() {
		if bev.Name == name {
			return bev
		}
	}
	return Beverage{}
}

// DefaultDrinks is the list of default drinks
// will be loading it from the DB
func DefaultDrinks() []Beverage {
	return []Beverage{
		Beverage{
			DefaultDrink:   true,
			Name:           "Double Espresso",
			Espresso:       espressoDouble,
			CupType:        cupSize6oz,
			Syrup:          syrupNone,
			Temperture:     tempHot,
			DrinkOfTheWeek: false,
		},
		Beverage{
			DefaultDrink:   true,
			Name:           "Hot Chocolate",
			Espresso:       espressoNone,
			CupType:        cupSize8oz,
			Syrup:          syrupNone,
			Temperture:     tempHot,
			DrinkOfTheWeek: false,
		},
		Beverage{
			DefaultDrink:   true,
			Name:           "Tea",
			Espresso:       espressoNone,
			CupType:        cupSize8oz,
			Syrup:          syrupHoney,
			Temperture:     tempHot,
			DrinkOfTheWeek: false,
		},
		Beverage{
			DefaultDrink:   true,
			Name:           "Iced Tea",
			Espresso:       espressoNone,
			CupType:        cupSize8oz,
			Syrup:          syrupHoney,
			Temperture:     tempIced,
			DrinkOfTheWeek: false,
		},
	}
}

// BeveragesForUser fetches the Beverages for that user
func BeveragesForUser(userID string) []Beverage {
	bevs := []Beverage{}
	db.Where(&Beverage{UserID: userID}).Find(&bevs)
	fmt.Println(bevs)
	return bevs
}
