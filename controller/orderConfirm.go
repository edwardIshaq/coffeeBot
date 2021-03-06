package controller

import (
	"SlackPlatform/middleware"
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

func (o *orderConfirmInteraction) postToStagingChannel(stagingChannelID string, beverage *models.Beverage, order models.Order, actionCallback SlackActionCallback, api *slack.Client) {
	callbackID := "order.confirmOrCancel"
	actionValue := fmt.Sprintf("%d", order.ID)
	postParams := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			slack.Attachment{
				CallbackID: callbackID,
				Text:       beverage.HumanReadable(),
				Actions: []slack.AttachmentAction{
					slack.AttachmentAction{
						Type:  "button",
						Text:  "Confirm Order",
						Name:  "confirm_beverage",
						Value: actionValue,
					},
					slack.AttachmentAction{
						Type:  "button",
						Text:  "Cancel Order",
						Name:  "cancel_beverage",
						Value: actionValue,
					},
				},
			},
		},
	}

	userText := fmt.Sprintf("<@%s|%s>", actionCallback.User.ID, actionCallback.User.Name)
	title := fmt.Sprintf("New Order from %s", userText)
	_, respTimestamp, err := api.PostMessage(stagingChannelID, title, postParams)
	if err != nil {
		fmt.Printf("Error posting to #cafe_requests %v", err)
		return
	}
	order.StagingMsgID = respTimestamp
	order.Save()
	fmt.Printf("respTimestamp = %s", respTimestamp)
}

func (o *orderConfirmInteraction) handleCallback(w http.ResponseWriter, r *http.Request, actionCallback SlackActionCallback) {
	action := actionCallback.Actions[0]
	order := models.OrderByID(action.Value)

	switch action.Name {
	case "confirm_beverage":
		handleConfirm(r, actionCallback, order)

	case "cancel_beverage":
		handleCancel(r, actionCallback, order)

	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}

func handleConfirm(r *http.Request, actionCallback SlackActionCallback, order models.Order) {
	// order.Confirm()
	order.IsConfirmed = true
	order.Save()

	updatedMessage := actionCallback.OriginalMessage
	if len(updatedMessage.Attachments) < 1 {
		fmt.Println("Has no attachments")
		return
	}

	newTitle := fmt.Sprintf("%s - Confirmed :white_check_mark:", updatedMessage.Text)
	attachment := updatedMessage.Attachments[0]
	attachment.Actions = []slack.AttachmentAction{}
	attachment.Color = "#00cc00"

	msgOption := slack.MsgOptionText(newTitle, false)
	attachmentOption := slack.MsgOptionAttachments(attachment)
	updateOption := slack.MsgOptionUpdate(updatedMessage.Timestamp)
	getSlackClientFromRequest(r)
	api.SendMessage(actionCallback.Channel.ID, updateOption, msgOption, attachmentOption)

	//Split into a seperate function
	//Post to #ProductionChannelID channel
	prodChannelID, ok := middleware.ProductionChannelID(r.Context())
	if !ok {
		fmt.Println("Failed to get `ProductionChannelID`")
		return
	}

	// add ready and cancel buttons
	readyButton := slack.AttachmentAction{
		Type: "button",
		Name: "order_ready",
		Text: "Ready to pickup",
	}
	cancelButton := slack.AttachmentAction{
		Type: "button",
		Name: "order_cancelled",
		Text: "Cancel order",
	}
	attachment.Actions = []slack.AttachmentAction{readyButton, cancelButton}
	attachment.CallbackID = "production_handler"
	attachmentOption = slack.MsgOptionAttachments(attachment)
	sendOption := slack.MsgOptionPost()

	_, respTimestamp, _, err := api.SendMessage(prodChannelID, sendOption, msgOption, attachmentOption)
	if err != nil {
		fmt.Printf("\nError posting to prod channel %v", err)
		return
	}

	order.ProdMsgID = respTimestamp
	order.Save()
}

func handleCancel(r *http.Request, actionCallback SlackActionCallback, order models.Order) {
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

	msgOption := slack.MsgOptionText(newTitle, false)
	attachmentOption := slack.MsgOptionAttachments(attachment)
	updateOption := slack.MsgOptionUpdate(updatedMessage.Timestamp)

	getSlackClientFromRequest(r)
	api.SendMessage(actionCallback.Channel.ID, updateOption, msgOption, attachmentOption)
}
