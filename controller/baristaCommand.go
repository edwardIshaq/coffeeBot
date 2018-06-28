package controller

/*
TODO:
	- Embedd the dialogInteractionHandler object instead of repeating the callbackID and regex
	- `handleCallback` is dedicated to one slash command and the class is made to look as a generic slash command ...
	- Check regex in dialogInteractionHandler can probably extract the ID
*/

import (
	"SlackPlatform/models"
	"fmt"
	"net/http"
	"regexp"

	"github.com/edwardIshaq/slack"
)

type slashCommand struct {
	// user will invoke `/<slashCommand>` to call this function
	slashCommand string
	// command will generate an interactive component with `callbackID` callback
	callbackID string
	//Regex to match against callbacks
	callbackRegex *regexp.Regexp
}

func baristaCommand() *slashCommand {
	callbackID := "beverage_selection"
	regex, _ := regexp.Compile(callbackID)
	return &slashCommand{
		slashCommand:  "coffeeCommand",
		callbackID:    callbackID,
		callbackRegex: regex,
	}
}

func (s *slashCommand) canHandleCallback(callback string) bool {
	return s.callbackRegex.MatchString(callback)
}

func (s *slashCommand) registerRoutes() {
	http.HandleFunc(s.route(), func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.sendBeverageMenu(w, r)
	})
}

func (s *slashCommand) route() string {
	return "/" + s.slashCommand
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
		CallbackID: s.callbackID,
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

func (s *slashCommand) handleCallback(w http.ResponseWriter, r *http.Request, actionCallback slack.AttachmentActionCallback) {
	_, ok := getSlackClientFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	r.ParseForm()

	if len(actionCallback.Actions) < 1 {
		return
	}
	if len(actionCallback.Actions[0].SelectedOptions) < 1 {
		return
	}

	selectedBevID := actionCallback.Actions[0].SelectedOptions[0].Value
	//Check if its a user defined beverage
	selectedBeverage := models.BeverageByID(selectedBevID)
	if !selectedBeverage.DefaultDrink {
		//Create a new Order
		order := models.SaveNewOrder(selectedBeverage)
		//if it is send a feedback message
		params := selectedBeverage.FeedbackMessage()
		replyMessage(params, actionCallback.ResponseURL)

		//Post to #cafeRequestsChannel
		go postToCafeChannel(&selectedBeverage, order, actionCallback, api)

		return
	}

	//udpate order with the triggerID so we can remove the menu message
	SlashBaristaMsgID := actionCallback.MessageTs
	triggerID := actionCallback.TriggerID
	if len(SlashBaristaMsgID) > 0 {
		order := models.OrderByBaristaMessageID(SlashBaristaMsgID)
		order.DialogTriggerID = triggerID
		order.BeverageID = selectedBeverage.ID
		fmt.Printf("order = %v", order)
		order.Save()
	}

	//Else picked one of the pre-defined bevs
	//Time to customize it
	dialog := selectedBeverage.MakeDialog(triggerID)
	api.OpenDialog(dialog)
}
