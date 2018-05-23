package models

import (
	"github.com/nlopes/slack"
)

// Dialog as in Slack dialogs
// 	https://api.slack.com/dialogs#option_element_attributes#top-level_dialog_attributes
type Dialog struct {
	TriggerID      string        `json:"trigger_id,omitempty"`
	CallbackID     string        `json:"callback_id"`
	NotifyOnCancel bool          `json:"notify_on_cancel"`
	Title          string        `json:"title"`
	SubmitLabel    string        `json:"submit_label,omitempty"`
	Elements       []interface{} `json:"elements"`
}

// DialogInput for dialogs input type text or menu
type DialogInput struct {
	Type        InputType `json:"type"`
	Label       string    `json:"label"`
	Name        string    `json:"name"`
	Placeholder string    `json:"placeholder"`
	Optional    bool      `json:"optional"`
}

// InputType is the type of the dialog input type
type InputType string

const (
	// InputTypeText textfield input
	InputTypeText InputType = "text"
	// InputTypeTextArea textarea input
	InputTypeTextArea InputType = "textarea"
	// InputTypeSelect textfield input
	InputTypeSelect InputType = "select"
)

// DialogTitle makes a title into a dialog title by caping it of to 24 chars
func DialogTitle(title string) string {
	const maxLength = 24
	if len(title) < maxLength {
		return title
	}
	return title[:21] + "..."
}

// DialogSubmitCallback to parse the response back from the Dialog
type DialogSubmitCallback struct {
	Type       string            `json:"type"`
	Submission map[string]string `json:"submission"`
	CallbackID string            `json:"callback_id"`

	Team        slack.Team    `json:"team"`
	Channel     slack.Channel `json:"channel"`
	User        slack.User    `json:"user"`
	ActionTs    string        `json:"action_ts"`
	Token       string        `json:"token"`
	ResponseURL string        `json:"response_url"`
}

// FeedbackMessage reply to dialog with an attachment message
func (d DialogSubmitCallback) FeedbackMessage(chosenBev string) *slack.Msg {
	fields := []slack.AttachmentField{}

	d.Submission["Name"] = chosenBev
	go func(d DialogSubmitCallback) {
		saveBeverage(d.Submission, d.User.ID)
	}(d)

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
				Text:   chosenBev,
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

// ExtractNames from the []Beverage
func ExtractNames(bevs []Beverage) []string {
	names := make([]string, 10)
	for _, bev := range bevs {
		names = append(names, bev.Name)
	}
	return names
}
