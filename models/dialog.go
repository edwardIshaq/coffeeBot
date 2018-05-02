package models

import (
	"encoding/json"
	"fmt"
)

// Dialog as in Slack dialogs
// 	https://api.slack.com/dialogs#option_element_attributes#top-level_dialog_attributes
type Dialog struct {
	TriggerID      string
	CallbackID     string             `json:"callback_id"`
	NotifyOnCancel bool               `json:"notify_on_cancel"`
	Title          string             `json:"title"`
	SubmitLabel    string             `json:"submit_label"`
	Elements       []baseInputElement `json:"elements"`
}

// ------------------------------------------
// DEMO
// ------------------------------------------

// DemoTextAreaElement this public
type DemoTextAreaElement struct {
	baseInputElement
	AdditionalField string `json:"moreFields"`
}

func notMain() {
	myObject := &DemoTextAreaElement{
		baseInputElement: baseInputElement{
			Type:  "textarea",
			Label: "My label",
			Name:  "Input name",
		},
		AdditionalField: "another thing",
	}

	if bts, err := json.Marshal(myObject); err == nil {
		fmt.Println(string(bts))
	}
}
