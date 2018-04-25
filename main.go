package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/nlopes/slack"
)

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
		<a href="https://slack.com/oauth/authorize?client_id=75950428352.351830378721&scope=commands,groups:read"><img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcset="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x" /></a>
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
	maybeCode, _ := exchangeCodeToAuth(code)
	message := "didnt work"
	if len(maybeCode) > 9 {
		message = fmt.Sprintf("worked last digits: %s", maybeCode[:10])
	}
	w.Write([]byte(message))
}

func redirectURL() string {
	return "https://" + appURL + "/oauthRedirect"
}

func exchangeCodeToAuth(code string) (string, error) {
	accessToken, scope, err := slack.GetOAuthToken(slackClientID, slackClientSecret, code, redirectURL(), false)
	fmt.Printf("accessToken: %s, scope: %s err: %s", accessToken, scope, err)
	return accessToken, err
}
