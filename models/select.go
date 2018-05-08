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

// baseSelect a menu select for dialogs
type baseSelect struct {
	DialogInput
	DataSource SelectDataSource `json:"data_source"`
}

//------------------------------------------
//		StaticSelectDialogInput
//------------------------------------------

// StaticSelectDialogInput can support all type except Dynamic menu
type StaticSelectDialogInput struct {
	baseSelect
	Value   string         `json:"value"` //This option is invalid in external, where you must use selected_options
	Options []SelectOption `json:"options"`
}

// SelectOption is an option for the user to select from the menu
type SelectOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func makeOptions(options []string) []SelectOption {
	selectOptions := make([]SelectOption, len(options))
	for idx, value := range options {
		selectOptions[idx] = SelectOption{
			Label: value,
			Value: value,
		}
	}
	return selectOptions
}

// NewStaticMenu constructor for a `static` datasource menu input
func NewStaticMenu(name, label string, options []string) *StaticSelectDialogInput {
	return &StaticSelectDialogInput{
		baseSelect: baseSelect{
			DialogInput: DialogInput{
				Type:  InputTypeSelect,
				Name:  name,
				Label: label,
			},
			DataSource: StaticDataSource,
		},
		Options: makeOptions(options),
	}
}

//------------------------------------------
//		GroupedSelectDialogInput
//------------------------------------------

// GroupedSelectDialogInput same as `StaticSelectDialogInput` but with grouped options
type GroupedSelectDialogInput struct {
	baseSelect
	Value        string
	OptionGroups []OptionGroup `json:"option_groups"`
}

// OptionGroup is a collection of options for creating a segmented table
type OptionGroup struct {
	Label   string         `json:"label"`
	Options []SelectOption `json:"options"`
}

// NewGroupedSelectDialoginput a grouped options select input for Dialogs
func NewGroupedSelectDialoginput(name, label string, groups map[string][]string) *GroupedSelectDialogInput {
	optionGroups := []OptionGroup{}
	for groupName, options := range groups {
		optionGroups = append(optionGroups, OptionGroup{
			Label:   groupName,
			Options: makeOptions(options),
		})
	}
	return &GroupedSelectDialogInput{
		baseSelect: baseSelect{
			DialogInput: DialogInput{
				Type:  InputTypeSelect,
				Name:  name,
				Label: label,
			},
			DataSource: StaticDataSource,
		},
		OptionGroups: optionGroups,
	}
}

//------------------------------------------
//		DynamicSelectInputElement
//------------------------------------------

// DynamicSelectInputElement special case for Dynamic since regular menu cant hold `value`
type DynamicSelectInputElement struct {
	baseSelect
}

//------------------------------------------
//		ExternalSelectInputElement
//------------------------------------------

// ExternalSelectInputElement is a special case of `SelectInputElement``
type ExternalSelectInputElement struct {
	baseSelect
	SelectedOptions []SelectOption `json:"selected_options"` //This option is invalid in external, where you must use selected_options
}
