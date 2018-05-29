package controller

import (
	"SlackPlatform/crossfunction"
	"SlackPlatform/middleware"
	"SlackPlatform/models"
	"net/http"
	"regexp"

	"github.com/nlopes/slack"
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

		api = crossfunction.ClientForRequest(r)
		s.respondToCommand(w, r)
	})
}

func (s *slashCommand) route() string {
	return "/" + s.slashCommand
}

func (s *slashCommand) respondToCommand(w http.ResponseWriter, r *http.Request) {
	channelID := r.PostFormValue("channel_id")
	attachment := slack.Attachment{}
	attachment.Text = "Choose a beverage"
	attachment.Fallback = "Choose a beverage from the menu"
	attachment.Color = "#3AA3E3"
	attachment.CallbackID = s.callbackID

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

	menuAction := slack.AttachmentAction{
		Name:         "beverage_menu",
		Text:         "Select beverage",
		Type:         "select",
		OptionGroups: optionGroups,
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

func (s *slashCommand) handleCallback(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	token, ok := middleware.AccessToken(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Process callback to extract `barista.dialog.<chosenBev>`
	actionCallback := parseAttachmentActionCallback(r)

	if len(actionCallback.Actions) < 1 {
		return
	}
	if len(actionCallback.Actions[0].SelectedOptions) < 1 {
		return
	}

	beverageSelectionID := actionCallback.Actions[0].SelectedOptions[0].Value
	//Check if its a user defined beverage
	presetBeverage := models.BeverageByID(beverageSelectionID)
	if !presetBeverage.DefaultDrink {
		//if it is send a feedback message
		params := presetBeverage.FeedbackMessage()
		replyMessage(params, actionCallback.ResponseURL)
		return
	}
	dialog := presetBeverage.MakeDialog()
	postDialog(dialog, actionCallback.TriggerID, token)
}
