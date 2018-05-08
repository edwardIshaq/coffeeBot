package models

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
