package models

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/jinzhu/gorm"
)

// Beverage type will hold the default selections for the menus
type Beverage struct {
	gorm.Model
	Category       string
	Name           string
	Espresso       string
	Syrup          string
	Temperture     string
	CupType        string
	Milk           string
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

	bev := BeverageByID(submission["ID"])
	submission["Name"] = bev.Name
	delete(submission, "ID")

	var err error
	for k, v := range submission {
		err = setField(&bev, k, v)
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

// BeverageByID finds a beverage by ID (string)
func BeverageByID(id string) Beverage {
	bevID, err := strconv.Atoi(id)
	if err != nil {
		return Beverage{}
	}
	beverage := Beverage{}
	db.First(&beverage, bevID)
	return beverage
}

// BeveragesForUser fetches the Beverages for that user
func BeveragesForUser(userID string) []Beverage {
	bevs := []Beverage{}
	db.Where(&Beverage{UserID: userID}).Find(&bevs)
	return bevs
}

// UserBeveragesMenu returns [pk]Name
func UserBeveragesMenu(userID string) map[string]string {
	list := BeveragesForUser(userID)
	menuMap := make(map[string]string)
	for _, bev := range list {
		menuMap[string(bev.ID)] = bev.Name
	}
	return menuMap
}

// MenuMap extract a map[bev.ID]bev.Name from a []Beverage
func MenuMap(bevs []Beverage) map[string]string {
	result := make(map[string]string)
	for _, bev := range bevs {
		result[fmt.Sprint(bev.ID)] = bev.Name
	}
	return result
}

// BeveragesByCategory select standard bevs by a specific category
func BeveragesByCategory(category string) []Beverage {
	result := []Beverage{}
	db.Where(&Beverage{Category: category, DefaultDrink: true}).Find(&result)
	return result
}
