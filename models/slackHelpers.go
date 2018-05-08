package models

import (
	"encoding/json"

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

// GetSelectedOption digs into the message menu response and gets the selected option
func GetSelectedOption(data []byte) (selectedOption, triggerID string) {
	var parsedResponse menuResponse
	selectedOption = ""
	triggerID = ""
	if err := json.Unmarshal(data, &parsedResponse); err != nil {
		return
	}

	triggerID = parsedResponse.TriggerID
	if len(parsedResponse.Actions) > 0 &&
		len(parsedResponse.Actions[0].SelectedOptions) > 0 {
		selectedOption = parsedResponse.Actions[0].SelectedOptions[0].Value
	}
	return
}
