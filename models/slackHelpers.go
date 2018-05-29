package models

import (
	"github.com/nlopes/slack"
)

// MakeAttachmentOptions converts an array of strings to slack menu options
func MakeAttachmentOptions(ss []string) []slack.AttachmentActionOption {
	actionOptions := []slack.AttachmentActionOption{}
	for _, value := range ss {
		actionOptions = append(actionOptions,
			slack.AttachmentActionOption{
				Text:  value,
				Value: value,
			})
	}
	return actionOptions
}

// MakeAttachmentOptionsFromMap converts a map of `map[string]string` to slack menu options
func MakeAttachmentOptionsFromMap(m map[string]string) []slack.AttachmentActionOption {
	actionOptions := []slack.AttachmentActionOption{}
	for value, text := range m {
		actionOptions = append(actionOptions,
			slack.AttachmentActionOption{
				Text:  text,
				Value: value,
			})
	}
	return actionOptions
}

// PayloadResponse a generic action response
type PayloadResponse struct {
	ResponseType      string     `json:"type"`
	CallbackID        string     `json:"callback_id"`
	Team              slack.Team `json:"team"`
	User              slack.User `json:"user"`
	VerificationToken string     `json:"token"`
}

// FeedbackMessage generates a slack feedback message for the chose beverage
func (b Beverage) FeedbackMessage() *slack.Msg {
	fields := &[]slack.AttachmentField{}
	appendFieldIfNotEmpty(fields, "Name", b.Name)
	appendFieldIfNotEmpty(fields, "Category", b.Category)
	appendFieldIfNotEmpty(fields, "Cup size", b.CupType)
	appendFieldIfNotEmpty(fields, "Espresso", b.Espresso)
	appendFieldIfNotEmpty(fields, "Syrup", b.Syrup)
	appendFieldIfNotEmpty(fields, "Temperture", b.Temperture)
	appendFieldIfNotEmpty(fields, "Comment", b.Comment)

	params := &slack.Msg{
		// Timestamp: d.ActionTs,
		Attachments: []slack.Attachment{
			slack.Attachment{
				// Text:   b.Name,
				Color:  "#eaca67",
				Fields: *fields,
			},
			slack.Attachment{
				Text:       "Confirm your here to pickup your order",
				CallbackID: "saveOrder",
				Actions: []slack.AttachmentAction{
					slack.AttachmentAction{
						Type:  "button",
						Text:  "Confirm",
						Name:  "confirm_order",
						Value: "confirm_order",
					},
				},
			},
		},
	}
	return params
}

func appendFieldIfNotEmpty(fields *[]slack.AttachmentField, title, value string) {
	if len(value) > 0 {
		*fields = append(*fields, slack.AttachmentField{
			Title: title,
			Value: value,
		})
	}
}

// SaveNewBeverage saves a new beverage from the dialog submission
func (d DialogSubmitCallback) SaveNewBeverage(chosenBevID string) *Beverage {
	templateBeverage := BeverageByID(chosenBevID)
	return saveBeverage(d.Submission, d.User.ID, templateBeverage)
}

// FeedbackMessage reply to dialog with an attachment message
func (d DialogSubmitCallback) FeedbackMessage(chosenBevID string) *slack.Msg {
	templateBeverage := BeverageByID(chosenBevID)

	fields := []slack.AttachmentField{}
	for key, value := range d.Submission {
		fields = append(fields,
			slack.AttachmentField{
				Title: key,
				Value: value,
			},
		)
	}

	params := &slack.Msg{
		Timestamp: d.ActionTs,
		Attachments: []slack.Attachment{
			slack.Attachment{
				Text:   templateBeverage.Name,
				Color:  "#eaca67",
				Fields: fields,
			},
			slack.Attachment{
				Text:       "Would you like to Save your Order for next time?",
				CallbackID: "saveOrder",
				Actions: []slack.AttachmentAction{
					slack.AttachmentAction{
						Type:  "button",
						Name:  "Save",
						Text:  "Save",
						Value: "SaveBeverage",
					},
				},
			},
		},
	}
	return params
}

// MakeDialog creates a new dialog from the `Beverage`
func (b Beverage) MakeDialog() Dialog {
	cupMenu := NewStaticSelectDialogInput("CupType", "Drink Size", AllDrinkSizes())
	cupMenu.Value = b.CupType

	espressMenu := NewStaticSelectDialogInput("Espresso", "Espresso Options", AllEspressoOptions())
	espressMenu.Value = b.Espresso

	syrupMenu := NewStaticSelectDialogInput("Syrup", "Syrup", AllSyrupOptions())
	syrupMenu.Value = b.Syrup

	tempMenu := NewStaticSelectDialogInput("Temperture", "Temperture", AllTemps())
	tempMenu.Value = b.Temperture

	commentInput := NewTextAreaInput("Comment", "Comments", b.Comment)
	commentInput.Optional = true

	callbackID := "barista.dialog." + string(b.ID)

	dialog := Dialog{
		CallbackID:  callbackID,
		Title:       DialogTitle(b.Name),
		SubmitLabel: "Order",
		Elements: []interface{}{
			cupMenu,
			espressMenu,
			syrupMenu,
			tempMenu,
			commentInput,
		},
	}
	return dialog
}
