package controller

import (
	"SlackPlatform/crossfunction"
	"SlackPlatform/models"
	"fmt"
	"net/http"
	"regexp"

	"github.com/edwardIshaq/slack"
)

type orderConfirmInteraction struct {
	// command will generate an interactive component with `callbackID` callback
	callbackID string
	//Regex to match against callbacks
	callbackRegex *regexp.Regexp
}

func newOrderConfirm() *orderConfirmInteraction {
	callbackID := "order.confirmOrCancel"
	// pattern := fmt.Sprintf(`%s\.(\d*)`, callbackID)
	regex, _ := regexp.Compile(`order\.confirmOrCancel`)
	return &orderConfirmInteraction{
		callbackID:    callbackID,
		callbackRegex: regex,
	}
}

func (o *orderConfirmInteraction) canHandleCallback(callback string) bool {
	return o.callbackRegex.MatchString(callback)
}

func (o *orderConfirmInteraction) callbackForID(beverageID uint) string {
	return fmt.Sprintf("%s.%d", o.callbackID, beverageID)
}

func (o *orderConfirmInteraction) handleCallback(w http.ResponseWriter, r *http.Request, actionCallback slack.AttachmentActionCallback) {
	// matches := o.callbackRegex.FindStringSubmatch(actionCallback.CallbackID)
	// if len(matches) < 2 {
	// 	fmt.Println("no matches", matches)
	// 	return
	// }

	action := actionCallback.Actions[0]
	order := models.OrderByID(action.Value)

	switch action.Name {
	case "confirm_beverage":
		handleConfirm(r, actionCallback, order)

	case "cancel_beverage":
		handleCancel(r, actionCallback, order)

	default:
		fmt.Println(action)
		w.WriteHeader(http.StatusBadRequest)
	}

}

func handleConfirm(r *http.Request, actionCallback slack.AttachmentActionCallback, order models.Order) {
	order.Confirm()
	updatedMessage := actionCallback.OriginalMessage
	if len(updatedMessage.Attachments) < 1 {
		fmt.Println("Has no attachments")
		return
	}

	newTitle := fmt.Sprintf("%s - Confirmed :white_check_mark:", updatedMessage.Text)
	attachment := updatedMessage.Attachments[0]
	attachment.Actions = []slack.AttachmentAction{}
	attachment.Color = "#00cc00"

	api = crossfunction.ClientForRequest(r)
	msgOption := slack.MsgOptionText(newTitle, false)
	attachmentOption := slack.MsgOptionAttachments(attachment)
	updateOption := slack.MsgOptionUpdate(updatedMessage.Timestamp)
	api.SendMessage(actionCallback.Channel.ID, updateOption, msgOption, attachmentOption)
}

func handleCancel(r *http.Request, actionCallback slack.AttachmentActionCallback, order models.Order) {
	order.Cancel()
	updatedMessage := actionCallback.OriginalMessage
	if len(updatedMessage.Attachments) < 1 {
		fmt.Println("Has no attachments")
		return
	}

	newTitle := fmt.Sprintf("%s - Canceled :no_entry:", updatedMessage.Text)

	attachment := updatedMessage.Attachments[0]
	attachment.Fields = []slack.AttachmentField{}
	attachment.Actions = []slack.AttachmentAction{}
	attachment.Color = "#BEBEBE"

	api = crossfunction.ClientForRequest(r)
	msgOption := slack.MsgOptionText(newTitle, false)
	attachmentOption := slack.MsgOptionAttachments(attachment)
	updateOption := slack.MsgOptionUpdate(updatedMessage.Timestamp)
	api.SendMessage(actionCallback.Channel.ID, updateOption, msgOption, attachmentOption)
}
