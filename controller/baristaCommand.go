package controller

import (
	"SlackPlatform/crossfunction"
	"SlackPlatform/models"
	"net/http"

	"github.com/nlopes/slack"
)

type slashCommand struct {
	name string
	api  *slack.Client
}

func (s *slashCommand) registerRoutes() {
	http.HandleFunc(s.route(), func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}
		s.api = crossfunction.ClientForRequest(r)
		s.handleCoffeeCommand(w, r)
	})
}

func (s *slashCommand) route() string {
	return "/" + s.name
}

func (s *slashCommand) handleCoffeeCommand(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	triggerID := r.PostFormValue("trigger_id")
	channelID := r.PostFormValue("channel_id")
	s.sendMenu(triggerID, channelID)
}

func (s *slashCommand) sendMenu(triggerID, channelID string) {
	attachment := slack.Attachment{}
	attachment.Text = "Choose a beverage"
	attachment.Fallback = "Choose a beverage from the menu"
	attachment.Color = "#3AA3E3"
	attachment.CallbackID = "beverage_selection"

	// Coffees
	coffeeGroup := slack.AttachmentActionOptionGroup{
		Text:    "Coffee",
		Options: models.MakeAttachmentOptions(models.AllCoffees()),
	}

	//drinkOfTheWeekGroup
	drinkOfTheWeekGroup := slack.AttachmentActionOptionGroup{
		Text:    "Drink of the Week",
		Options: models.MakeAttachmentOptions(models.AllDrinksOfTheWeek()),
	}

	// Regular Drinks Menu
	regularDrinksGroup := slack.AttachmentActionOptionGroup{
		Text:    "Usual Drinks",
		Options: models.MakeAttachmentOptions(models.AllUsualDrinks()),
	}

	// Tea
	teaDrinksGroup := slack.AttachmentActionOptionGroup{
		Text:    "Tea",
		Options: models.MakeAttachmentOptions(models.AllTeas()),
	}

	menuAction := slack.AttachmentAction{
		Name:         "beverage_menu",
		Text:         "Select beverage",
		Type:         "select",
		OptionGroups: []slack.AttachmentActionOptionGroup{coffeeGroup, regularDrinksGroup, teaDrinksGroup, drinkOfTheWeekGroup},
	}

	attachment.Actions = []slack.AttachmentAction{menuAction}
	message := slack.Message{}
	message.Text = "What would you like to order ?"
	message.Attachments = []slack.Attachment{attachment}

	postParams := slack.NewPostMessageParameters()
	postParams.Attachments = []slack.Attachment{attachment}
	postParams.Channel = channelID

	s.api.PostMessage(channelID, "choose a beverage", postParams)
}
