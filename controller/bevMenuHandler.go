package controller

import (
	"SlackPlatform/models"
	"fmt"
	"net/http"
	"regexp"

	"github.com/edwardIshaq/slack"
)

type bevMenuHandler struct {
	// messageMenu will generate an interactive component with `callbackID` callback
	callbackID string
	//Regex to match against callbacks
	callbackRegex *regexp.Regexp
}

func newBevMenuHndler() *bevMenuHandler {
	callbackID := "beverage_selection"
	regex, _ := regexp.Compile(callbackID)
	return &bevMenuHandler{
		callbackID:    callbackID,
		callbackRegex: regex,
	}
}
func (b *bevMenuHandler) canHandleCallback(callback string) bool {
	return b.callbackRegex.MatchString(callback)
}

func (b *bevMenuHandler) handleCallback(w http.ResponseWriter, r *http.Request, actionCallback slack.AttachmentActionCallback) {
	_, ok := getSlackClientFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	r.ParseForm()

	if len(actionCallback.Actions) < 1 {
		return
	}
	if len(actionCallback.Actions[0].SelectedOptions) < 1 {
		return
	}

	selectedBevID := actionCallback.Actions[0].SelectedOptions[0].Value
	//Check if its a user defined beverage
	selectedBeverage := models.BeverageByID(selectedBevID)
	if !selectedBeverage.DefaultDrink {
		//Create a new Order
		order := models.SaveNewOrder(selectedBeverage)
		//if it is send a feedback message
		params := selectedBeverage.FeedbackMessage()
		replyMessage(params, actionCallback.ResponseURL)

		//Post to #cafeRequestsChannel
		go postToCafeChannel(&selectedBeverage, order, actionCallback, api)

		return
	}

	//udpate order with the triggerID so we can remove the menu message
	SlashBaristaMsgID := actionCallback.MessageTs
	triggerID := actionCallback.TriggerID
	if len(SlashBaristaMsgID) > 0 {
		order := models.OrderByBaristaMessageID(SlashBaristaMsgID)
		order.DialogTriggerID = triggerID
		order.BeverageID = selectedBeverage.ID
		fmt.Printf("order = %v", order)
		order.Save()
	}

	//Else picked one of the pre-defined bevs
	//Time to customize it
	dialog := selectedBeverage.MakeDialog(triggerID)
	api.OpenDialog(dialog)
}
