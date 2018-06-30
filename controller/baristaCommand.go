package controller

import (
	"SlackPlatform/models"
	"fmt"
	"net/http"

	"github.com/edwardIshaq/slack"
)

type slashCommand struct {
	// user will invoke `/<slashCommand>` to call this function
	slashCommand string
}

func baristaCommand() *slashCommand {
	return &slashCommand{"coffeeCommand"}
}

func (s *slashCommand) registerRoutes() {
	slashCommandRoute := "/" + s.slashCommand
	http.HandleFunc(slashCommandRoute, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.sendBeverageMenu(w, r)
	})
}

func (s *slashCommand) sendBeverageMenu(w http.ResponseWriter, r *http.Request) {
	_, ok := getSlackClientFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	channelID := r.PostFormValue("channel_id")
	attachment := slack.Attachment{
		Fallback:   "Choose a beverage from the menu",
		Color:      "#3AA3E3",
		CallbackID: beverageMenuHandler.callbackID,
	}

	//User Beverages
	userID := r.PostFormValue("user_id")
	userDrinks := slack.AttachmentActionOptionGroup{
		Text:    "User Beverages",
		Options: menuFromBevs(models.BeveragesForUser(userID)),
	}

	optionGroups := []slack.AttachmentActionOptionGroup{userDrinks}

	//The rest of the drinks from the DB
	allBevs := models.AllBeveragesByCategory()
	for category, beverages := range allBevs {
		optionGroups = append(optionGroups,
			slack.AttachmentActionOptionGroup{
				Text:    category,
				Options: menuFromBevs(beverages),
			})
	}

	attachment.Actions = []slack.AttachmentAction{
		slack.AttachmentAction{
			Name:         "beverage_menu",
			Text:         "Select beverage",
			Type:         "select",
			OptionGroups: optionGroups,
		},
	}

	attachmentOption := slack.MsgOptionAttachments(attachment)
	textOption := slack.MsgOptionText("What would you like to order ?", false)
	ephemeralOption := slack.MsgOptionPostEphemeral2(userID)

	baristaMessageID, err := api.PostEphemeral(channelID, userID, textOption, ephemeralOption, attachmentOption)
	if err != nil {
		fmt.Printf("Error sending slash command %v", err)
		return
	} else if baristaMessageID == "" {
		fmt.Println("No messageID received")
		return
	}

	//else Start a new order with the `slashMessageTS`
	go func(baristaMessageID string) {
		models.NewBaristaCommandOrder(baristaMessageID)
	}(baristaMessageID)
}

func menuFromBevs(bevs []models.Beverage) []slack.AttachmentActionOption {
	bevsMap := models.MenuMap(bevs)
	return models.MakeAttachmentOptionsFromMap(bevsMap)
}
