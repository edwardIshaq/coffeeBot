package models

// Dialog as in Slack dialogs
// 	https://api.slack.com/dialogs#option_element_attributes#top-level_dialog_attributes
type Dialog struct {
	TriggerID      string
	CallbackID     string             `json:"callback_id"`
	NotifyOnCancel bool               `json:"notify_on_cancel"`
	Title          string             `json:"title"`
	SubmitLabel    string             `json:"submit_label"`
	Elements       []TextInputElement `json:"elements"`
}
