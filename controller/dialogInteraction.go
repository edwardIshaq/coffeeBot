package controller

import (
	"SlackPlatform/middleware"
	"SlackPlatform/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/edwardIshaq/slack"
)

type dialogInteraction struct {
	// command will generate an interactive component with `callbackID` callback
	callbackID string
	//Regex to match against callbacks
	callbackRegex *regexp.Regexp
}

func beverageDialogInteraction() *dialogInteraction {
	callbackID := "barista.dialog"
	pattern := `barista\.dialog\.(.*)`
	regex, _ := regexp.Compile(pattern)

	return &dialogInteraction{
		callbackID:    callbackID,
		callbackRegex: regex,
	}
}

func (d *dialogInteraction) canHandleCallback(callback string) bool {
	return d.callbackRegex.MatchString(callback)
}

func (d *dialogInteraction) handleCallback(w http.ResponseWriter, r *http.Request, actionCallback SlackActionCallback) {
	matches := d.callbackRegex.FindStringSubmatch(actionCallback.CallbackID)
	if len(matches) < 2 {
		fmt.Printf("no triggerID found on: %s", actionCallback.CallbackID)
		return
	}

	//Fetch the correct order
	triggerID := matches[1]
	orderQuery := models.Order{DialogTriggerID: triggerID}
	fetchedOrder := orderQuery.Fetch()
	if fetchedOrder == nil {
		fmt.Printf("Couldn't find a matching order %v", actionCallback)
		return
	}

	// save beverage
	dialogResponse := slack.DialogSubmitCallback{}
	json.Unmarshal([]byte(r.PostFormValue("payload")), &dialogResponse)
	bevID := fmt.Sprintf("%d", fetchedOrder.BeverageID)
	beverage := models.SaveNewBeverage(dialogResponse, bevID)
	if beverage == nil {
		log.Println("Failed to save a new beverage")
		return
	}
	fetchedOrder.BeverageID = beverage.ID
	fetchedOrder.Save()
	orderID := strconv.FormatInt(int64(fetchedOrder.ID), 10)

	//post feedback message to user
	feedback := beverage.FeedbackMessage()
	buttonAttachment := slack.Attachment{
		CallbackID: saveBevAction.callbackID,
		Actions: []slack.AttachmentAction{
			//Save Beverage
			slack.AttachmentAction{
				Name:  "SaveButton",
				Text:  "Save",
				Type:  "button",
				Value: bevID,
			},
			//Cancel Order
			slack.AttachmentAction{
				Name:  "CancelButton",
				Text:  "Cancel",
				Type:  "button",
				Value: orderID,
			},
		},
	}
	attachments := feedback.Attachments
	attachments = append(attachments, buttonAttachment)
	feedback.Attachments = attachments
	replyMessage(feedback, actionCallback.ResponseURL)

	w.WriteHeader(http.StatusOK)

	//Post to #stagingChannelID
	stagingChannelID, ok := middleware.StagingChannelID(r.Context())
	if !ok {
		fmt.Println("Failed to get `stagingChannelID`")
		return
	}
	go orderConfirmHandler.postToStagingChannel(stagingChannelID, beverage, *fetchedOrder, actionCallback, api)

	channelID := actionCallback.Channel.ID
	fmt.Printf("Now delete the menu message: %s %s", fetchedOrder.SlashBaristaMsgID, channelID)

	str1, str2, err := api.DeleteMessage(channelID, fetchedOrder.SlashBaristaMsgID)
	fmt.Printf("%s | %s | %v", str1, str2, err)

	// go func(slashBaristaMsgID, channel string, api *slack.Client) {
	// 	str1, str2, err := api.DeleteMessage(slashBaristaMsgID, channelID)
	// 	fmt.Printf("%s | %s | %v", str1, str2, err)
	// }(fetchedOrder.SlashBaristaMsgID, channelID, api)
}
