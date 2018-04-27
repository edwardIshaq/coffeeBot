package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

/*
This will support a basic play grounds for the `utaApp` found at the link below.string
utaApp https://api.slack.com/apps/AABQEB4M7
The server is configured to run on:
https://goplatform.ngrok.io/

add to slack button:	https://goplatform.ngrok.io/addToSlack
OAuth handlers:			https://goplatform.ngrok.io/oauthRedirect
outgoing hook url: 		https://goplatform.ngrok.io/outgoingHooks

Notes:
------------------------------------------
Permissions update:
	When ever the permissions are changed on the App's scopes
	the button details (found in `func buttonTemplate()``) should be updated from https://api.slack.com/apps/AABQEB4M7/distribute?

TODO:
------------------------------------------
[] save tokens to DB

*/

const (
	utaAppToken = "xoxp-75950428352-75957863573-352376775186-cd1dd2d44890733760faaa6eda878916"
	appURL      = "goplatform.ngrok.io"

	slackClientID     = "75950428352.351830378721"
	slackClientSecret = "a56df86a6f1fae41f4efceaf20fb9842"
)

var (
	message = "Hello world"
)

func main() {
	//getGroups()
	http.HandleFunc("/addToSlack", installWTAApplication)
	http.HandleFunc("/oauthRedirect", oAuthRedirectHandler)
	http.HandleFunc("/outgoingHooks", handleOutgoingHooks)
	http.HandleFunc("/", sayHello)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

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

func buttonTemplate() string {
	button := `
	<body>
		Add this to Slack <br>
		<a href="https://slack.com/oauth/authorize?client_id=75950428352.351830378721&scope=commands,groups:read,bot,chat:write:bot">
		<img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcset="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x" />
		</a>
	</body>
	`
	return button
}

func appSuggestionURL() string {
	return `<meta name="slack-app-id" content="AABQEB4M7">`
}

func appSuggestionHTML() string {
	return `<meta name="slack-app-id" content="AABQEB4M7">`
}

func installWTAApplication(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(buttonTemplate()))
}

// code=75950428352.351913834208.ab4422d8cec8e7b134b8dbe7659e097cecb294122cfebee8456cad92f03d1732&state=
func oAuthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	//use `code` to get `OAuthResponse` response back
	oauthResponse, err := slack.GetOAuthResponseContext(context.Background(), slackClientID, slackClientSecret, code, redirectURL(), false)

	message := "Something wrong happened"
	if err == nil {
		message = ""
		message += fmt.Sprintf("\nError?: %v", err)
		message += fmt.Sprintf("\nBot token: %v", oauthResponse.Bot)
		message += fmt.Sprintf("\naccessToken: %v", oauthResponse.AccessToken)
		message += fmt.Sprintf("\nscopes: %v", oauthResponse.Scope)
	}
	w.Write([]byte(message))
}

func redirectURL() string {
	return "https://" + appURL + "/oauthRedirect"
}

//------------------------------------------
//Slack outgoing hooks demo
func handleOutgoingHooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new outgoingWebhook from Slack")
	r.ParseForm()
	fmt.Println(r.PostForm)
	text := r.PostFormValue("text")
	value, _ := strconv.ParseInt(text, 10, 32)
	message := fmt.Sprintf(`{"text" : "%d"}`, value+1)
	w.Write([]byte(message))
}
