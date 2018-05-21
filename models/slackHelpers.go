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
