package main

/*
This will support a basic play grounds for the `utaApp` found at the link below.string
utaApp https://api.slack.com/apps/AABQEB4M7
The server is configured to run on:
https://goplatform.ngrok.io/
https://slackcoffeebar.typeform.com/to/v6kODV

add to slack button:	https://goplatform.ngrok.io/addToSlack
OAuth handlers:			https://goplatform.ngrok.io/oauthRedirect
outgoing hook url: 		https://goplatform.ngrok.io/outgoingHooks
/barista command url:	https://goplatform.ngrok.io/coffeeCommand
interactive url:		https://goplatform.ngrok.io/interactive
Dynamic Menu:			https://goplatform.ngrok.io/dynamicMenu

TODO:
------------------------------------------
[] save tokens to DB
*/

import (
	"SlackPlatform/controller"
	"SlackPlatform/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/nlopes/slack"
)

const (
	appURL            = "goplatform.ngrok.io"
	slackClientID     = "75950428352.351830378721"
	slackClientSecret = "a56df86a6f1fae41f4efceaf20fb9842"
	verificationToken = "8ycguzKPPcWvt7wIsud0a9EL"
)

/*
utaApp info https://api.slack.com/apps/AABQEB4M7
info which is related to a single team `eddie-beta`
Needs to be stored in DB
*/
const (
	utaAppToken = "xoxp-75950428352-75957863573-355080493893-c39a5f8e88a4b08e475dbce0d0b4884e"
)

var (
	message = "Hello world"
)

func main() {
	db := connectToDatabase()
	defer db.Close()
	models.SetDatabase(db)
	//models.DemoDB()

	dbWrapper := models.NewDBWrapper(db)
	controller.StartupControllers(dbWrapper)

	//getGroups()
	http.HandleFunc("/outgoingHooks", handleOutgoingHooks)
	http.HandleFunc("/coffeeCommand", handleCoffeeCommand)
	http.HandleFunc("/interactive", handleInteractiveMessages)
	http.HandleFunc("/", sayHello)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func connectToDatabase() *sql.DB {
	connStr := "user=goUser dbname=barista password=qwe123 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

/*
Testing the api
*/
func getGroups() {
	api := slack.New(utaAppToken)

	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// api.SetDebug(true)
	groups, err := api.GetGroups(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	for _, group := range groups {
		var isPrivate = "public-group"
		if group.IsPrivate {
			isPrivate = "Private"
		}
		fmt.Printf("ID: %s \t Name: %s \tPrivate: %s\n", group.ID, group.Name, isPrivate)
	}
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	message += " " + r.Method
	message += "\nBody:" + fmt.Sprintf("%s", r.Body)
	w.Write([]byte(message))
}

//------------------------------------------
//Slack outgoing hooks demo
func handleOutgoingHooks(w http.ResponseWriter, r *http.Request) {
	text := r.PostFormValue("text")
	value, _ := strconv.ParseInt(text, 10, 32)
	message := fmt.Sprintf(`{"text" : "%d"}`, value+1)
	w.Write([]byte(message))
}

func handleCoffeeCommand(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	triggerID := r.PostFormValue("trigger_id")
	channelID := r.PostFormValue("channel_id")
	sendMenu(triggerID, channelID)
}

func sendMenu(triggerID, channelID string) {
	attachment := slack.Attachment{}
	attachment.Text = "Choose a beverage"
	attachment.Fallback = "hmmmm, something wrong"
	attachment.Color = "#3AA3E3"
	attachment.CallbackID = "beverage_selection"

	//drinkOfTheWeekGroup
	drinkOfTheWeekGroup := slack.AttachmentActionOptionGroup{
		Text:    "Drink of the Week",
		Options: models.MakeAttachmentOptions([]string{"Vitality Latte", "Herbal Remedy Tea", "Iced Separator"}),
	}

	// Regular Drinks Menu
	regularDrinksGroup := slack.AttachmentActionOptionGroup{
		Text:    "Usual Drinks",
		Options: models.MakeAttachmentOptions([]string{"Steamed Milk", "Hot Chocolate", "Tea"}),
	}

	teaDrinksGroup := slack.AttachmentActionOptionGroup{
		Text:    "Tea",
		Options: models.MakeAttachmentOptions([]string{"London Fog", "San-Fran Fog", "Matcha Latte", "Tanglewood Ginger Chai"}),
	}

	menuAction := slack.AttachmentAction{
		Name:         "beverage_menu",
		Text:         "Select beverage",
		Type:         "select",
		OptionGroups: []slack.AttachmentActionOptionGroup{regularDrinksGroup, teaDrinksGroup, drinkOfTheWeekGroup},
	}

	attachment.Actions = []slack.AttachmentAction{menuAction}
	message := slack.Message{}
	message.Text = "What would you like to order ?"
	message.Attachments = []slack.Attachment{attachment}

	postParams := slack.NewPostMessageParameters()
	postParams.Attachments = []slack.Attachment{attachment}
	postParams.Channel = channelID

	fmt.Print(postParams)
	api := slack.New(utaAppToken)
	api.PostMessage(channelID, "choose a beverage", postParams)
}

func readFile(filePath string) (content string, err error) {
	data, err := ioutil.ReadFile(filePath)
	content = string(data)
	return
}

func handleInteractiveMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n\nInteractive message")
	r.ParseForm()
	payloadString := r.PostFormValue("payload")
	fmt.Print(payloadString)
	fmt.Println("")
	res := payload{}
	json.Unmarshal([]byte(payloadString), &res)
	fmt.Printf("\n triggerID : %s", res.TriggerID)
}

type payload struct {
	TriggerID string `json:"trigger_id"`
}

func sendDialog() {

}
