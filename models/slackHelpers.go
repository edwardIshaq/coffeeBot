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

// ------------------------------------------
type menuResponse struct {
	Type      string   `json:"type"`
	Actions   []action `json:"actions"`
	TriggerID string   `json:"trigger_id"`
}

type action struct {
	Name            string        `json:"name"`
	Type            string        `json:"type"`
	SelectedOptions []optionValue `json:"selected_options"`
}

type optionValue struct {
	Value string `json:"value"`
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
						Name:  "confirm_order_name",
						Text:  "Confirm",
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
