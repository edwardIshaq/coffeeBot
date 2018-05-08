package models

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

// TextInputElement subtype of DialogInput
//	https://api.slack.com/dialogs#option_element_attributes#text_element_attributes
type TextInputElement struct {
	DialogInput
	MaxLength int              `json:"max_length,omitempty"`
	MinLength int              `json:"min_length,omitempty"`
	Hint      string           `json:"hint,omitempty"`
	Subtype   TextInputSubtype `json:"subtype"`
	Value     string           `json:"value"`
}

// NewTextInput constructor for a `text` input
func NewTextInput(name, label string) *TextInputElement {
	return &TextInputElement{
		DialogInput: DialogInput{
			Type:  InputTypeText,
			Name:  name,
			Label: label,
		},
	}
}

// NewTextAreaInput constructor for a `textarea` input
func NewTextAreaInput(name, label string) *TextInputElement {
	return &TextInputElement{
		DialogInput: DialogInput{
			Type:  InputTypeTextArea,
			Name:  name,
			Label: label,
		},
	}
}
