package controller

import (
	"SlackPlatform/middleware"
	"SlackPlatform/models"
	"fmt"
	"net/http"
	"regexp"
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

func (b *bevMenuHandler) handleCallback(w http.ResponseWriter, r *http.Request, actionCallback SlackActionCallback) {
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

	SlashBaristaMsgID := actionCallback.MessageTs
	selectedBevID := actionCallback.Actions[0].SelectedOptions[0].Value
	selectedBeverage := models.BeverageByID(selectedBevID)

	order := models.OrderByBaristaMessageID(SlashBaristaMsgID)
	order.BeverageID = selectedBeverage.ID

	//Check if its a user defined beverage
	if !selectedBeverage.DefaultDrink {
		//its a customized drink -> send a feedback message
		params := selectedBeverage.FeedbackMessage()
		replyMessage(params, actionCallback.ResponseURL)

		//Post to #cafeRequestsChannel
		stagingChannelID, ok := middleware.StagingChannelID(r.Context())
		if !ok {
			fmt.Println("Failed to get `stagingChannelID`")
			return
		}
		go postToStagingChannel(stagingChannelID, &selectedBeverage, order, actionCallback, api)
		return
	}

	//udpate order with the triggerID so we can remove the menu message
	triggerID := actionCallback.TriggerID
	if len(SlashBaristaMsgID) > 0 {
		order.DialogTriggerID = triggerID
		fmt.Printf("order = %v", order)
		order.Save()
	}

	//Else picked a `DefaultDrink`
	//Time to customize it
	//FIXME: Move to `dialogInteraction.go`
	dialog := selectedBeverage.MakeDialog(triggerID)
	api.OpenDialog(dialog)
}
