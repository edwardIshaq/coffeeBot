package models

// SelectDataSource types of select datasource
type SelectDataSource string

const (
	// StaticDataSource menu with static Options/OptionGroups
	StaticDataSource SelectDataSource = "static"
	// ExternalDataSource dynamic datasource
	ExternalDataSource SelectDataSource = "external"
	// ConversationsDataSource provides a list of conversations
	ConversationsDataSource SelectDataSource = "conversations"
	// ChannelsDataSource provides a list of channels
	ChannelsDataSource SelectDataSource = "channels"
	// UsersDataSource provides a list of users
	UsersDataSource SelectDataSource = "users"
)

// baseSelectInputElement a menu select for dialogs
type baseSelectInputElement struct {
	InputElement
	DataSource SelectDataSource `json:"data_source"`
}

// SelectInputElement can support all type except Dynamic menu
type SelectInputElement struct {
	baseSelectInputElement
	Value   string         `json:"value"` //This option is invalid in external, where you must use selected_options
	Options []SelectOption `json:"options"`
}

// NewStaticMenu constructor for a `static` datasource menu input
func NewStaticMenu(name, label string, options []string) *SelectInputElement {
	selectOptions := convertStringsToSelectOptions(options)
	return &SelectInputElement{
		baseSelectInputElement: baseSelectInputElement{
			InputElement: InputElement{
				Type:  InputTypeSelect,
				Name:  name,
				Label: label,
			},
			DataSource: StaticDataSource,
		},
		Options: selectOptions,
	}
}
func convertStringsToSelectOptions(options []string) []SelectOption {
	selectOptions := make([]SelectOption, len(options))
	for idx, value := range options {
		selectOptions[idx] = newSelectOption(value)
	}
	return selectOptions
}

// DynamicSelectInputElement special case for Dynamic since regular menu cant hold `value`
type DynamicSelectInputElement struct {
	baseSelectInputElement
}

// ExternalSelectInputElement is a special case of `SelectInputElement``
type ExternalSelectInputElement struct {
	baseSelectInputElement
	SelectedOptions Options `json:"selected_options"` //This option is invalid in external, where you must use selected_options
}

// Options an alias for `[]SelectOption`
type Options []SelectOption

// SelectOption is an option for the user to select from the menu
type SelectOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// newSelectOption will create an option with the `value` provided
func newSelectOption(value string) SelectOption {
	return SelectOption{
		Label: value,
		Value: value,
	}
}

// OptionGroup is a collection of options for creating a segmented table
type OptionGroup struct {
	Label   string         `json:"label"`
	Options []SelectOption `json:"options"`
}
