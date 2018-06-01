package controller

import (
	"SlackPlatform/middleware"
	"SlackPlatform/models"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

/*
Todo:
Rename the interactions to somethig shorter
*/
type saveBeverageInteraction struct {
	// command will generate an interactive component with `callbackID` callback
	callbackID string
	//Regex to match against callbacks
	callbackRegex *regexp.Regexp
}

func saveBevInteraction() *saveBeverageInteraction {
	callbackID := "order_created"
	regex, _ := regexp.Compile(callbackID)

	return &saveBeverageInteraction{
		callbackID:    callbackID,
		callbackRegex: regex,
	}
}

func (d *saveBeverageInteraction) canHandleCallback(callback string) bool {
	return d.callbackRegex.MatchString(callback)
}

func (d *saveBeverageInteraction) handleCallback(w http.ResponseWriter, r *http.Request) {
	token, ok := middleware.AccessToken(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	actionCallback := parseAttachmentActionCallback(r)
	fmt.Printf("\n actionCallback = %v \n", actionCallback)

	// saveActionValue := fmt.Sprintf("save_beverage.%d", b.ID)

	if len(actionCallback.Actions) < 1 {
		fmt.Println("There is no Actions, the impossible happened")
		return
	}

	actionComponents := strings.Split(actionCallback.Actions[0].Value, ".")
	if len(actionComponents) < 1 {
		fmt.Println("There is no BevID")
		return
	}
	chosenBevID := actionComponents[1]

	presetBeverage := models.BeverageByID(chosenBevID)
	dialog := presetBeverage.MakeSaveNameDialog()
	postDialog(dialog, actionCallback.TriggerID, token)
}
