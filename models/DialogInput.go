package models

// DialogInput for dialogs input type text or menu
type DialogInput struct {
	Type        InputType `json:"type"`
	Label       string    `json:"label"`
	Name        string    `json:"name"`
	Placeholder string    `json:"placeholder"`
	Optional    bool      `json:"optional"`
}
