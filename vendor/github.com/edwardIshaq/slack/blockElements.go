package slack

/*
TODO:
[ ] enforce limit on Overflow 10 options limit
[ ] enforce limit on Select 10 options limit

*/

// BlockElementType type of elements that could be contained in a Block
type BlockElementType string

const (
	// SelectBlockElementType a select element
	SelectBlockElementType BlockElementType = "select"
	// ButtonBlockElementType a button element
	ButtonBlockElementType BlockElementType = "button"
	// DatePickerBlockElementType a datePicker element
	DatePickerBlockElementType BlockElementType = "datepicker"
	// OverflowBlockElementType an overflow menu
	OverflowBlockElementType BlockElementType = "overflow"
	// ImageBlockElementType an image block element
	ImageBlockElementType BlockElementType = "image"
	// TextBlockElementType simple text block
	TextBlockElementType BlockElementType = "text"
	// UserBlockElementType a user block
	UserBlockElementType BlockElementType = "user"
)

// BlockElement generic BlockElement type
type BlockElement interface{}

// TextBlockCompatibleElement an interface to check if a BlockElement can be embedded in `TextBlock`
type TextBlockCompatibleElement interface {
	canEmbeddInTextBlock() bool
}

// TextBlockElement a simple text block element
type TextBlockElement struct {
	Type BlockElementType `json:"type"`
	Text string           `json:"text"` //Required
}

// NewTextBlockElement convenience func to create `TextBlockElement`
func NewTextBlockElement(text string) TextBlockElement {
	return TextBlockElement{
		Type: TextBlockElementType,
		Text: text,
	}
}

// ButtonBlockElement a button action block element
type ButtonBlockElement struct {
	Type     BlockElementType `json:"type"`
	Text     string           `json:"text"` //Required
	ActionID string           `json:"action_id"`
	Value    string           `json:"value"`
}

func (b ButtonBlockElement) canEmbeddInTextBlock() bool {
	return true
}

// NewButtonBlockElement convnience method to create `ButtonBlockElement`
func NewButtonBlockElement(text, actionID, value string) ButtonBlockElement {
	return ButtonBlockElement{
		Type:     ButtonBlockElementType,
		Text:     text,
		ActionID: actionID,
		Value:    value,
	}
}

// OverflowBlockElement a small context menu of up to 10 options
type OverflowBlockElement struct {
	Type     BlockElementType    `json:"type"`
	ActionID string              `json:"action_id"`
	Options  []SelectBlockOption `json:"options"` //Limited to 10
}

func (o OverflowBlockElement) canEmbeddInTextBlock() bool {
	return true
}

// NewOverflowBlockElement convnience method to create `OverflowBlockElement`
func NewOverflowBlockElement(actionID string, options []SelectBlockOption) OverflowBlockElement {
	return OverflowBlockElement{
		Type:     OverflowBlockElementType,
		ActionID: actionID,
		Options:  options,
	}
}

// SelectBlockDataSource types of select datasource
type SelectBlockDataSource string

const (
	// StaticSelectDataSource menu with static Options/OptionGroups
	StaticSelectDataSource SelectBlockDataSource = "static"
	// ExternalSelectDataSource dynamic datasource
	ExternalSelectDataSource SelectBlockDataSource = "external"
	// ConversationsSelectDataSource provides a list of conversations
	ConversationsSelectDataSource SelectBlockDataSource = "conversations"
	// ChannelsSelectDataSource provides a list of channels
	ChannelsSelectDataSource SelectBlockDataSource = "channels"
	// UsersSelectDataSource provides a list of users
	UsersSelectDataSource SelectBlockDataSource = "users"
	// DateSelectDataSource a date picker data source
	DateSelectDataSource SelectBlockDataSource = "date"
)

// SelectBlockElement is a `Select Element` of type `select`
type SelectBlockElement struct {
	Type        BlockElementType      `json:"type"`
	ActionID    string                `json:"action_id,omitempty"`
	PlaceHolder string                `json:"placeholder,omitempty"`
	DataSource  SelectBlockDataSource `json:"data_source,omitempty"`
	Options     []SelectBlockOption   `json:"options,omitempty"`
}

func (s SelectBlockElement) canEmbeddInTextBlock() bool {
	return true
}

// NewSelectBlockElement convenience Constructor to create a `SelectBlockElement`
func NewSelectBlockElement(actionID, placeHolder string, dataSource SelectBlockDataSource, options []SelectBlockOption) SelectBlockElement {
	return SelectBlockElement{
		Type:        SelectBlockElementType,
		ActionID:    actionID,
		PlaceHolder: placeHolder,
		DataSource:  dataSource,
		Options:     options,
	}
}

// SelectBlockOption the individual option to appear in a `Select` Block element
type SelectBlockOption struct {
	Text  string `json:"text"`  // Required.
	Value string `json:"value"` // Required.
}

// ImageBlockElement an image block
type ImageBlockElement struct {
	Type     BlockElementType `json:"type"`
	ImageURL string           `json:"image_url,omitempty"`
	AltText  string           `json:"alt_text,omitempty"`
}

func (i ImageBlockElement) canEmbeddInTextBlock() bool {
	return true
}

// NewImageBlockElemenet convnience method to create `ImageBlockElement`
func NewImageBlockElemenet(imageURL, altText string) ImageBlockElement {
	return ImageBlockElement{
		Type:     ImageBlockElementType,
		ImageURL: imageURL,
		AltText:  altText,
	}
}

// DatePickerBlockElement a datepicker element type
type DatePickerBlockElement struct {
	Type     BlockElementType `json:"type"`
	ActionID string           `json:"action_id,omitempty"`
	Value    string           `json:"value,omitempty"`
}

func (d DatePickerBlockElement) canEmbeddInTextBlock() bool {
	return true
}

// NewDatePickerBlockElement convnience method to create `DatePickerBlockElement`
func NewDatePickerBlockElement(actionID, value string) DatePickerBlockElement {
	return DatePickerBlockElement{
		Type:     DatePickerBlockElementType,
		ActionID: actionID,
		Value:    value,
	}
}
