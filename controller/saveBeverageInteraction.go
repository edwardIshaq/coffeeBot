package controller

import (
	"SlackPlatform/middleware"
	"SlackPlatform/models"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/edwardIshaq/slack"
)

/*
Todo:
[ ] Rename the interactions to somethig shorter
[ ] Display Drink Name in the title
[ ] Show a confirmation after saving
*/
type saveBeverageInteraction struct {
	// command will generate an interactive component with `callbackID` callback
	callbackID string
	//Regex to match against callbacks
	callbackRegex *regexp.Regexp
}

func saveBevInteraction() *saveBeverageInteraction {
	callbackID := "order_created"
	regex, _ := regexp.Compile(`order_created|saveBeverageName\.(\d*)`)
	return &saveBeverageInteraction{
		callbackID:    callbackID,
		callbackRegex: regex,
	}
}

func (d *saveBeverageInteraction) canHandleCallback(callback string) bool {
	return d.callbackRegex.MatchString(callback)
}

func (d *saveBeverageInteraction) handleCallback(w http.ResponseWriter, r *http.Request, actionCallback slack.AttachmentActionCallback) {
	token, ok := middleware.AccessToken(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if len(actionCallback.Actions) < 1 {
		dialogResponse := slack.DialogSubmitCallback{}
		json.Unmarshal([]byte(r.PostFormValue("payload")), &dialogResponse)
		if len(dialogResponse.Submission) > 0 {
			d.handleSaveNameCallback(w, r, dialogResponse)
			return
		}
		return
	}

	actionComponents := strings.Split(actionCallback.Actions[0].Value, ".")
	if len(actionComponents) < 2 {
		fmt.Printf("Expecting `<action>.<bevID>` format, instead got %s\n", actionCallback.Actions[0].Value)
		return
	}
	action := actionComponents[0]
	chosenBevID := actionComponents[1]
	beverage := models.BeverageByID(chosenBevID)

	switch action {
	case "save_beverage":
		dialog := beverage.MakeSaveNameDialog()
		postDialog(dialog, actionCallback.TriggerID, token)

	case "confirm_beverage":
		fmt.Printf("Confirm order %v", beverage)

	}

}

func (d *saveBeverageInteraction) handleSaveNameCallback(w http.ResponseWriter,
	r *http.Request,
	submitCallback slack.DialogSubmitCallback) {
	matches := d.callbackRegex.FindStringSubmatch(submitCallback.CallbackID)
	if len(matches) < 2 {
		fmt.Println("no matches", matches)
		return
	}
	beverageID := matches[1]
	bev := models.BeverageByID(beverageID)
	drinkName := submitCallback.Submission["drinkName"]
	comment := submitCallback.Submission["comment"]
	bev.UpdateDrink(drinkName, comment)
	w.WriteHeader(http.StatusOK)
}
