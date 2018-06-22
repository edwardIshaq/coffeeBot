package controller

import (
	"SlackPlatform/middleware"
	"SlackPlatform/models"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

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
	callbackID := "NameBeverageOrCancelOrder"
	pattern := `NameBeverageOrCancelOrder|NameBeverage\.(\d*)`
	regex, _ := regexp.Compile(pattern)
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

	switch actionCallback.CallbackID {
	case "NameBeverageOrCancelOrder":
		d.handleButtonPressed(r, actionCallback, token)

	default:
		//NameBeverage.<Beverage.ID>
		if len(actionCallback.Actions) == 0 {
			dialogResponse := slack.DialogSubmitCallback{}
			json.Unmarshal([]byte(r.PostFormValue("payload")), &dialogResponse)
			if len(dialogResponse.Submission) > 0 {
				d.handleSaveNameCallback(w, r, dialogResponse, actionCallback)
				return
			}
			return
		}
	}

}

func (d *saveBeverageInteraction) handleButtonPressed(r *http.Request, actionCallback slack.AttachmentActionCallback, token string) {
	if len(actionCallback.Actions) != 1 {
		fmt.Printf("\nExpecting 1 Action got %v", actionCallback.Actions)
		return
	}

	action := actionCallback.Actions[0]
	switch action.Name {
	case "SaveButton":
		beverage := models.BeverageByID(action.Value)
		dialog := beverage.MakeSaveNameDialog()
		postDialog(dialog, actionCallback.TriggerID, token)

	case "CancelButton":
		orderID := action.Value
		order := models.OrderByID(orderID)
		handleCancel(r, actionCallback, order)
	}
}

func (d *saveBeverageInteraction) handleSaveNameCallback(w http.ResponseWriter,
	r *http.Request,
	submitCallback slack.DialogSubmitCallback,
	actionCallback slack.AttachmentActionCallback) {

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

	params := &slack.Msg{
		Text: fmt.Sprintf("New drink saved as %s", drinkName),
	}
	replyMessage(params, actionCallback.ResponseURL)
}
