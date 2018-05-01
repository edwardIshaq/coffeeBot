package models

// Dialog as in Slack dialogs
// 	https://api.slack.com/dialogs#option_element_attributes#top-level_dialog_attributes
type Dialog struct {
	TriggerID      string
	CallbackID     string         `json:"callback_id"`
	NotifyOnCancel bool           `json:"notify_on_cancel"`
	Title          string         `json:"title"`
	SubmitLabel    string         `json:"submit_label"`
	Elements       []InputElement `json:"elements"`
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

// TextInputSubtype Accepts email, number, tel, or url. In some form factors, optimized input is provided for this subtype.
type TextInputSubtype string

const (
	// EmailTextInputSubtype email keyboard
	EmailTextInputSubtype TextInputSubtype = "email"
	// NumberTextInputSubtype numeric keyboard
	NumberTextInputSubtype TextInputSubtype = "number"
	// TelTextInputSubtype Phone keyboard
	TelTextInputSubtype TextInputSubtype = "tel"
	// URLTextInputSubtype Phone keyboard
	URLTextInputSubtype TextInputSubtype = "url"
)

// SelectDataSource types of select datasource
type SelectDataSource string

const (
	// ExternalDataSource dynamic datasource
	ExternalDataSource SelectDataSource = "external"
	// ConversationsDataSource provides a list of conversations
	ConversationsDataSource SelectDataSource = "conversations"
	// ChannelsDataSource provides a list of channels
	ChannelsDataSource SelectDataSource = "channels"
	// UsersDataSource provides a list of users
	UsersDataSource SelectDataSource = "users"
)

// InputElement for dialogs input type text or menu
type InputElement struct {
	// Type of form element. For a text input, the type is always text. Required.
	Type InputType `json:"type"`
	// Label displayed to user. Required. 24 character maximum.
	Label string `json:"label"`
	// Name of form element. Required. No more than 300 characters
	Name string `json:"name"`
}

// TextInputElement subtype of InputElement
//	https://api.slack.com/dialogs#option_element_attributes#text_element_attributes
type TextInputElement struct {
	Type        InputType        `json:"type"`
	Label       string           `json:"label"`
	Name        string           `json:"name"`
	MaxLength   int              `json:"max_length"`
	MinLength   int              `json:"min_length"`
	Optional    bool             `json:"optional"`
	Hint        string           `json:"hint"`
	Subtype     TextInputSubtype `json:"subtype"`
	Value       string           `json:"value"`
	Placeholder string           `json:"placeholder"`
}
