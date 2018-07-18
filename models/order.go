package models

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

// Order tracks the status of a Beverage order
type Order struct {
	gorm.Model

	//Slack flow
	SlashBaristaMsgID string
	DialogTriggerID   string
	BeverageID        uint
	StagingMsgID      string
	ProdMsgID         string

	IsConfirmed bool
	IsFulfilled bool
}

// NewBaristaCommandOrder Creates an order in the DB with the /barista messageID
func NewBaristaCommandOrder(baristaMessageID string) *Order {
	order := Order{
		SlashBaristaMsgID: baristaMessageID,
	}
	db.Create(&order)
	return &order
}

// OrderByBaristaMessageID finds a beverage by BaristaCommandMessageID
func OrderByBaristaMessageID(id string) Order {
	order := Order{}
	db.Where(&Order{SlashBaristaMsgID: id}).First(&order)
	return order
}

// Fetch fetch a specified order
func (o Order) Fetch() *Order {
	fetchedOrder := new(Order)
	db.Where(&o).First(&fetchedOrder)
	return fetchedOrder
}

// SaveNewOrder Creates an order in the DB
func SaveNewOrder(b Beverage) Order {
	order := Order{
		BeverageID:  b.ID,
		IsConfirmed: false,
		IsFulfilled: false,
	}
	db.Create(&order)
	return order
}

// OrderByID finds a beverage by ID (string)
func OrderByID(id string) Order {
	orderID, err := strconv.Atoi(id)
	if err != nil {
		return Order{}
	}
	order := Order{}
	db.First(&order, orderID)
	return order
}

// Confirm update the model to be confirmed
// FIXME: Not working properly !!!
func (o Order) Confirm() {
	o.IsConfirmed = true
	db.Save(o)
}

// Cancel update the model to be confirmed
func (o Order) Cancel() {
	db.Delete(o)
}

// Save saves the order
func (o Order) Save() {
	db.Save(o)
}
