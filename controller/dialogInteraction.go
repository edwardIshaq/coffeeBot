package controller

import (
	"SlackPlatform/models"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type dialogInteraction struct {
	// command will generate an interactive component with `callbackID` callback
	callbackID string
	//Regex to match against callbacks
	callbackRegex *regexp.Regexp
}

func beverageDialogInteraction() *dialogInteraction {
	callbackID := "barista.dialog"
	regex, _ := regexp.Compile(callbackID)

	return &dialogInteraction{
		callbackID:    callbackID,
		callbackRegex: regex,
	}
}

func (d *dialogInteraction) canHandleCallback(callback string) bool {
	return d.callbackRegex.MatchString(callback)
}

func (d *dialogInteraction) handleCallback(w http.ResponseWriter, r *http.Request) {
	actionCallback := parseAttachmentActionCallback(r)
	var chosenBev string
	if strings.HasPrefix(actionCallback.CallbackID, "barista.dialog.") {
		chosenBev = strings.Split(actionCallback.CallbackID, ".")[2]
	}

	dialogResponse := models.DialogSubmitCallback{}
	json.Unmarshal([]byte(r.PostFormValue("payload")), &dialogResponse)
	// save beverage and post feedback message
	newBeverage := dialogResponse.SaveNewBeverage(chosenBev)
	if newBeverage == nil {
		log.Println("Failed to save a new beverage")
		return
	}
	params := newBeverage.FeedbackMessage()
	replyMessage(params, actionCallback.ResponseURL)
}
