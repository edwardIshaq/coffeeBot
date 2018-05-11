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
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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
	api     = slack.New(utaAppToken)
)

func main() {
	db := connectToDatabase()
	defer db.Close()
	models.SetDatabase(db)
	//models.DemoDB()

	dbWrapper := models.NewDBWrapper(db)
	controller.StartupControllers(dbWrapper, api)

	http.HandleFunc("/outgoingHooks", handleOutgoingHooks)
	// http.HandleFunc("/coffeeCommand", handleCoffeeCommand)
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

func parseAttachmentActionCallback(r *http.Request) slack.AttachmentActionCallback {
	r.ParseForm()
	payload := r.PostFormValue("payload")
	actionCallback := slack.AttachmentActionCallback{}
	json.Unmarshal([]byte(payload), &actionCallback)
	return actionCallback
}

func handleInteractiveMessages(w http.ResponseWriter, r *http.Request) {
	actionCallback := parseAttachmentActionCallback(r)
	fmt.Println()
	callbackID := actionCallback.CallbackID
	chosenBev := ""
	if strings.HasPrefix(actionCallback.CallbackID, "barista.dialog.") {
		callbackID = "barista.dialog"
		chosenBev = strings.Split(actionCallback.CallbackID, ".")[2]
	}

	switch callbackID {
	case "beverage_selection":
		fmt.Println("interacted with `menu`")
		//textReply(w, "Customize your order")
		if len(actionCallback.Actions) >= 1 {
			if len(actionCallback.Actions[0].SelectedOptions) >= 1 {
				chosenBeverage := actionCallback.Actions[0].SelectedOptions[0].Value
				fmt.Printf("You selected %v", chosenBeverage)
				fmt.Printf("triggerID: %v", actionCallback.TriggerID)
				postDialog(chosenBeverage, actionCallback.TriggerID)
			}
		}
		return

	case "barista.dialog":
		startTime := time.Now()
		dialogResponse := models.DialogSubmitCallback{}
		json.Unmarshal([]byte(r.PostFormValue("payload")), &dialogResponse)
		responseURL := dialogResponse.ResponseURL
		params := dialogResponse.FeedbackMessage(chosenBev)

		go func(params *slack.Msg, responseURL string) {
			data, _ := json.Marshal(params)
			bodyReader := bytes.NewReader(data)
			req, err := http.NewRequest(http.MethodPost, responseURL, bodyReader)
			fmt.Println("Request\n", req)

			//Fire the request
			resp, err := slack.HTTPClient.Do(req)
			if err != nil {
				fmt.Println("\nResponseError: ", err)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("RESPONSE: %v", resp)
			fmt.Printf("\nprocessing Dialog took: %s\n", time.Since(startTime))
		}(params, responseURL)
	}
}

func textReply(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "application/json")
	message := fmt.Sprintf(`{"text" : "%s", "replace_original": false}`, text)
	w.Write([]byte(message))
}

func postDialog(chosenBeverage, triggerID string) {
	dialog := makeDialog(chosenBeverage)

	if dialogjson, err := json.Marshal(dialog); err == nil {
		postBody := url.Values{
			"token":      {"xoxp-75950428352-75957863573-355080493893-c39a5f8e88a4b08e475dbce0d0b4884e"},
			"trigger_id": {triggerID},
			"dialog":     {string(dialogjson)},
		}

		reqBody := strings.NewReader(postBody.Encode())
		req, err := http.NewRequest("POST", slack.SLACK_API+"dialog.open", reqBody)
		if err != nil {
			fmt.Println("error happened: ", err)
			return
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = req.WithContext(context.Background())

		//Fire the request
		resp, err := slack.HTTPClient.Do(req)
		if err != nil {
			fmt.Println("\nResponseError: ", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				fmt.Printf("error reading body %v", err2)
				return
			}
			bodyString := string(bodyBytes)
			fmt.Println("\nbodyString: ", bodyString)
		}
	}
}

func makeDialog(chosenBeverage string) models.Dialog {
	presetBeverage := models.BeverageByName(chosenBeverage)

	cupMenu := models.NewStaticSelectDialogInput("CupType", "Drink Size", models.AllDrinkSizes())
	cupMenu.Value = presetBeverage.CupType

	espressMenu := models.NewStaticSelectDialogInput("Espresso", "Espresso Options", models.AllEspressoOptions())
	espressMenu.Value = presetBeverage.Espresso

	syrupMenu := models.NewStaticSelectDialogInput("Syrup", "Syrup", models.AllSyrupOptions())
	syrupMenu.Value = presetBeverage.Syrup

	tempMenu := models.NewStaticSelectDialogInput("Temperture", "Temperture", models.AllTemps())
	tempMenu.Value = presetBeverage.Temperture

	commentInput := models.NewTextAreaInput("Comment", "Comments")
	commentInput.Optional = true

	callbackID := "barista.dialog." + chosenBeverage

	dialog := models.Dialog{
		CallbackID:  callbackID,
		Title:       models.DialogTitle(chosenBeverage),
		SubmitLabel: "Order",
		Elements: []interface{}{
			cupMenu,
			espressMenu,
			syrupMenu,
			tempMenu,
			commentInput,
		},
	}
	return dialog
}
