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

// DefaultDrinks is the list of default drinks
// will be loading it from the DB
func DefaultDrinks() []Beverage {
	return []Beverage{
		Beverage{
			DefaultDrink: true,
			Name:         "Espresso",
			Espresso:     espressoSingle,
		},
		Beverage{
			DefaultDrink: true,
			Name:         "Hot Chocolate",
			Espresso:     espressoNone,
		},
	}
}
