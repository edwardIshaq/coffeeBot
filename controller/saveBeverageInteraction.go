package controller

import (
	"SlackPlatform/middleware"
	"SlackPlatform/models"
	"encoding/json"
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
	regex, _ := regexp.Compile(`order_created|saveBeverageName\.(\d*)`)
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
	if len(actionCallback.Actions) < 1 {
		dialogResponse := models.DialogSubmitCallback{}
		json.Unmarshal([]byte(r.PostFormValue("payload")), &dialogResponse)
		if len(dialogResponse.Submission) > 0 {
			d.handleSaveNameCallback(w, r, dialogResponse)
			return
		}
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

func (d *saveBeverageInteraction) handleSaveNameCallback(w http.ResponseWriter,
	r *http.Request,
	submitCallback models.DialogSubmitCallback) {
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
