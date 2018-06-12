package models

import (
	"fmt"

	"github.com/edwardIshaq/slack"
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

// CreateFields converts a beverage to a list of fields
func (b Beverage) CreateFields() []slack.AttachmentField {
	fields := &[]slack.AttachmentField{}
	appendFieldIfNotEmpty(fields, "Name", b.Name)
	appendFieldIfNotEmpty(fields, "Category", b.Category)
	appendFieldIfNotEmpty(fields, "Cup size", b.CupType)
	appendFieldIfNotEmpty(fields, "Espresso", b.Espresso)
	appendFieldIfNotEmpty(fields, "Syrup", b.Syrup)
	appendFieldIfNotEmpty(fields, "Temperture", b.Temperture)
	appendFieldIfNotEmpty(fields, "Comment", b.Comment)
	return *fields
}

// FeedbackMessage generates a slack feedback message for the chose beverage
func (b Beverage) FeedbackMessage() *slack.Msg {
	fields := b.CreateFields()

	saveActionValue := fmt.Sprintf("save_beverage.%d", b.ID)
	params := &slack.Msg{
		Attachments: []slack.Attachment{
			slack.Attachment{
				Text:   "Please confirm your order in channel #cafe_requests when you arrive",
				Color:  "#eaca67",
				Fields: fields,
			},
			slack.Attachment{
				Text:       "Would you like to name this drink for future orders?",
				CallbackID: "order_created",
				Actions: []slack.AttachmentAction{
					slack.AttachmentAction{
						Type:  "button",
						Text:  "Save",
						Name:  "save_beverage",
						Value: saveActionValue,
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
func SaveNewBeverage(d slack.DialogSubmitCallback, chosenBevID string) *Beverage {
	templateBeverage := BeverageByID(chosenBevID)
	return saveBeverage(d.Submission, d.User.ID, templateBeverage)
}

// MakeDialog creates a new dialog from the `Beverage`
func (b Beverage) MakeDialog() slack.Dialog {
	cupMenu := slack.NewStaticSelectDialogInput("CupType", "Drink Size", stringsToSelectOptions(AllDrinkSizes()))
	cupMenu.Value = b.CupType

	espressMenu := slack.NewStaticSelectDialogInput("Espresso", "Espresso Options", stringsToSelectOptions(AllEspressoOptions()))
	espressMenu.Value = b.Espresso

	syrupMenu := slack.NewStaticSelectDialogInput("Syrup", "Syrup", stringsToSelectOptions(AllSyrupOptions()))
	syrupMenu.Value = b.Syrup

	tempMenu := slack.NewStaticSelectDialogInput("Temperture", "Temperture", stringsToSelectOptions(AllTemps()))
	tempMenu.Value = b.Temperture

	milkMenu := slack.NewStaticSelectDialogInput("Milk", "Milk", stringsToSelectOptions(AllMilkOptions()))
	milkMenu.Value = b.Milk

	callbackID := fmt.Sprintf("barista.dialog.%d", b.ID)

	dialog := slack.Dialog{
		CallbackID:  callbackID,
		Title:       DialogTitle(b.Name),
		SubmitLabel: "Order",
		Elements: []slack.DialogElement{
			cupMenu,
			espressMenu,
			syrupMenu,
			tempMenu,
			milkMenu,
		},
	}
	return dialog
}

type stringArray []string

func stringsToSelectOptions(options []string) []slack.SelectOption {
	selectOptions := make([]slack.SelectOption, len(options))
	for idx, value := range options {
		selectOptions[idx] = slack.SelectOption{
			Label: value,
			Value: value,
		}
	}
	return selectOptions
}

// DialogTitle makes a title into a dialog title by caping it of to 24 chars
func DialogTitle(title string) string {
	const maxLength = 24
	if len(title) < maxLength {
		return title
	}
	return title[:21] + "..."
}

// MakeSaveNameDialog to save a custom title for the drink
func (b Beverage) MakeSaveNameDialog() slack.Dialog {
	callbackID := fmt.Sprintf("saveBeverageName.%d", b.ID)

	nameInput := slack.NewTextInput("drinkName", "Drink Name", b.Name)
	commentInput := slack.NewTextAreaInput("comment", "Comments", b.Comment)
	commentInput.Optional = true

	return slack.Dialog{
		CallbackID:  callbackID,
		Title:       fmt.Sprintf("Save Drink"),
		SubmitLabel: "Save",
		Elements: []slack.DialogElement{
			nameInput,
			commentInput,
		},
	}
}
