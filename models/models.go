package models

import (
	"github.com/nlopes/slack"
)

// MenuOption for Dialogs or messageMenu
type MenuOption struct {
	Title string `json:"text"`
	Value string `json:"value"`
}

// NewMenuOption creates a `Menu` with same title/value
func NewMenuOption(title string) MenuOption {
	return MenuOption{
		Title: title,
		Value: title,
	}
}

func (m MenuOption) attachmentOption() slack.AttachmentActionOption {
	return slack.AttachmentActionOption{
		Text:  m.Title,
		Value: m.Value,
	}
}

// Menu is an array of `Menu`
type Menu []MenuOption

// ConvertToAttachmentActionOptions converts Menus ->
func (ms Menu) ConvertToAttachmentActionOptions() []slack.AttachmentActionOption {
	var options []slack.AttachmentActionOption

	for _, m := range ms {
		option := m.attachmentOption()
		options = append(options, option)
	}
	return options
}

func convertToMenuOptions(ss []string) Menu {
	var menuOptions []MenuOption
	for _, s := range ss {
		menuOptions = append(menuOptions, NewMenuOption(s))
	}
	return menuOptions
}

// MakeAttachmentOptions converts an array of strings to slack menu options
func MakeAttachmentOptions(ss []string) []slack.AttachmentActionOption {
	menu := convertToMenuOptions(ss)
	return menu.ConvertToAttachmentActionOptions()
}
