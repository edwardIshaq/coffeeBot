package controller

import (
	"SlackPlatform/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/edwardIshaq/slack"
)

const (
	cafeRequestsChannel = "C0SCMST6C"
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

func (d *dialogInteraction) handleCallback(w http.ResponseWriter, r *http.Request, actionCallback slack.AttachmentActionCallback) {
	var chosenBev string
	if strings.HasPrefix(actionCallback.CallbackID, "barista.dialog.") {
		chosenBev = strings.Split(actionCallback.CallbackID, ".")[2]
	}

	dialogResponse := slack.DialogSubmitCallback{}
	json.Unmarshal([]byte(r.PostFormValue("payload")), &dialogResponse)

	// save beverage
	beverage := models.SaveNewBeverage(dialogResponse, chosenBev)
	if beverage == nil {
		log.Println("Failed to save a new beverage")
		return
	}

	//post feedback message to user
	replyMessage(beverage.FeedbackMessage(), actionCallback.ResponseURL)

	//Post to #cafeRequestsChannel
	go postToCafeChannel(beverage, actionCallback, api)
}

func postToCafeChannel(beverage *models.Beverage, actionCallback slack.AttachmentActionCallback, api *slack.Client) {
	confirmActionValue := fmt.Sprintf("confirm_beverage.%d", beverage.ID)
	postParams := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			slack.Attachment{
				CallbackID: "order.confirmOrCancel",
				Fields:     beverage.CreateFields(),
				Actions: []slack.AttachmentAction{
					slack.AttachmentAction{
						Type:  "button",
						Text:  "Confirm Order",
						Name:  "confirm_beverage",
						Value: confirmActionValue,
					},
					slack.AttachmentAction{
						Type:  "button",
						Text:  "Cancel Order",
						Name:  "cancel_beverage",
						Value: confirmActionValue,
					},
				},
			},
		},
	}
	title := fmt.Sprintf("New Order from *%s*", actionCallback.User.Name)
	if _, _, err := api.PostMessage(cafeRequestsChannel, title, postParams); err != nil {
		fmt.Printf("Error posting to #cafe_requests %v", err)
	}
}