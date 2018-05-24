package controller

import (
	"SlackPlatform/crossfunction"
	"SlackPlatform/models"
	"net/http"

	"github.com/nlopes/slack"
)

type slashCommand struct {
	name string
}

func (s *slashCommand) registerRoutes() {
	http.HandleFunc(s.route(), func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		api = crossfunction.ClientForRequest(r)
		s.handleCoffeeCommand(w, r)
	})
}

func (s *slashCommand) route() string {
	return "/" + s.name
}

func (s *slashCommand) handleCoffeeCommand(w http.ResponseWriter, r *http.Request) {
	channelID := r.PostFormValue("channel_id")
	attachment := slack.Attachment{}
	attachment.Text = "Choose a beverage"
	attachment.Fallback = "Choose a beverage from the menu"
	attachment.Color = "#3AA3E3"
	attachment.CallbackID = "beverage_selection"

	//User Beverages
	userID := r.PostFormValue("user_id")
	userDrinks := slack.AttachmentActionOptionGroup{
		Text:    "User Beverages",
		Options: menuFromBevs(models.BeveragesForUser(userID)),
	}

	// Coffees
	coffeeGroup := slack.AttachmentActionOptionGroup{
		Text:    "Coffee",
		Options: menuFromBevs(models.BeveragesByCategory("Coffee")),
	}

	//drinkOfTheWeekGroup
	drinkOfTheWeekGroup := slack.AttachmentActionOptionGroup{
		Text:    "Drink of the week",
		Options: menuFromBevs(models.BeveragesByCategory("Drink of the week")),
	}

	// Tea
	teaDrinksGroup := slack.AttachmentActionOptionGroup{
		Text:    "Tea",
		Options: menuFromBevs(models.BeveragesByCategory("Tea")),
	}

	menuAction := slack.AttachmentAction{
		Name:         "beverage_menu",
		Text:         "Select beverage",
		Type:         "select",
		OptionGroups: []slack.AttachmentActionOptionGroup{userDrinks, coffeeGroup, teaDrinksGroup, drinkOfTheWeekGroup},
	}

	attachment.Actions = []slack.AttachmentAction{menuAction}
	message := slack.Message{}
	message.Text = "What would you like to order ?"
	message.Attachments = []slack.Attachment{attachment}

	postParams := slack.NewPostMessageParameters()
	postParams.Attachments = []slack.Attachment{attachment}
	postParams.Channel = channelID

	api.PostMessage(channelID, "choose a beverage", postParams)
}

func menuFromBevs(bevs []models.Beverage) []slack.AttachmentActionOption {
	bevsMap := models.MenuMap(bevs)
	return models.MakeAttachmentOptionsFromMap(bevsMap)
}
