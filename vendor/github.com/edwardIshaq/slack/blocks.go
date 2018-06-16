package slack

// BlockType captures the type of block
type BlockType string

const (
	// DividerBlockType horizontal divider
	DividerBlockType BlockType = "divider"
	// FileBlockType a file block
	FileBlockType BlockType = "file"
	// ImageBlockType an image block
	ImageBlockType BlockType = "image"
	// TextBlockType a simple text block
	TextBlockType BlockType = "text"
	// TextCollectionBlockType Text Collection Block
	TextCollectionBlockType BlockType = "text_collection"
	// ActionBlockType is a collection of up to 5 Actions (select, button, overflow, datepicker)
	ActionBlockType BlockType = "action"
	// ContextBlockType is a collection of up to 10 elements of type (image, text, user)
	ContextBlockType BlockType = "context"
	// CanvasBlockType a canvas (TBD)
	CanvasBlockType BlockType = "canvas"
)

// TextBlockSubType element types that fit in `TextBlockType`
type TextBlockSubType string

const (
	// ImageBlockSubType embedded image in text
	ImageBlockSubType TextBlockSubType = "image"
	// ButtonBlockSubType embedded button in text
	ButtonBlockSubType TextBlockSubType = "button"
)

// Block interface for the attachments
type Block interface{}

// BaseBlock base block element
type BaseBlock struct {
	BlockID string `json:"block_id,omitempty"`
	Type    string `json:"type"`
}

// DividerBlock a horizontal devider block
type DividerBlock struct {
	BaseBlock
}

// NewDividerBlock builds a devider block
func NewDividerBlock(blockID string) DividerBlock {
	return DividerBlock{
		BaseBlock: BaseBlock{
			BlockID: blockID,
			Type:    string(DividerBlockType),
		},
	}
}

// TextBlock a simple text block
type TextBlock struct {
	BaseBlock
	Text    string       `json:"text,omitempty"`
	Element BlockElement `json:"element,omitempty"`
}

// NewTextBlock constructs a new text block with ID and text
func NewTextBlock(text, blockID string) TextBlock {
	return TextBlock{
		BaseBlock: BaseBlock{
			BlockID: blockID,
			Type:    string(TextBlockType),
		},
		Text: text,
	}
}

// ActionBlock base block element
type ActionBlock struct {
	ActionID string `json:"action_id,omitempty"`
	Type     string `json:"type"`
}

/*
{
	"type": "image",
	"block_id": "image4",
	"image_url": "https://scontent-sjc3-1.cdninstagram.com/vp/64d7aa4ab1a55892036c52b2237f3868/5B97948F/t51.2885-15/s640x640/sh0.08/e35/c0.135.1080.1080/26155058_127776921353331_485398838513762304_n.jpg",
	"alt_text": "cat",
	"caption": "Bubsy eats a cheeto"
}
*/

// ImageBlock a text block with an image embedded in it
type ImageBlock struct {
	BaseBlock
	ImageURL string `json:"image_url,omitempty"`
	AltText  string `json:"alt_text,omitempty"`
}

// NewImageBlock constructs a new text block with ID and text
func NewImageBlock(imageURL, altText, blockID string) ImageBlock {
	return ImageBlock{
		BaseBlock: BaseBlock{
			BlockID: blockID,
			Type:    string(ImageBlockSubType),
		},
		ImageURL: imageURL,
		AltText:  altText,
	}
}
