package slack

// BlockElement generic BlockElement type
type BlockElement interface{}

// BaseBlockElement basic building blocks
type BaseBlockElement struct {
	Type string `json:"type,omitempty"`
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
	BaseBlockElement
	ActionID    string                `json:"action_id,omitempty"`
	PlaceHolder string                `json:"placeholder,omitempty"`
	DataSource  SelectDataSource      `json:"data_source,omitempty"`
	Options     [10]SelectBlockOption `json:"options,omitempty"`
}

// SelectBlockOption the individual option to appear in a `Select` Block element
type SelectBlockOption struct {
	Text  string `json:"text"`  // Required.
	Value string `json:"value"` // Required.
}

// ButtonBlockElement a button block element
type ButtonBlockElement struct {
	BaseBlockElement
	Text     string `json:"text"` // Required.
	ActionID string `json:"action_id,omitempty"`
	Value    string `json:"value,omitempty"`
}
