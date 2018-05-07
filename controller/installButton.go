package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nlopes/slack"
)

const (
	defaultScopes = "commands,groups:read,bot,chat:write:bot"
)

type appInstaller struct {
	slack.OAuthResponse
	appURL            string
	slackClientID     string
	slackClientSecret string
	verificationToken string
}

/*
Handler for the install button
*/
func (installer *appInstaller) registerRoutes() {
	http.HandleFunc("/addToSlack", installer.installWTAApplication)
	http.HandleFunc("/oauthRedirect", installer.oAuthRedirectHandler)
}

// InstallWTAApplication registers a route to install the app on workspaces
func (installer *appInstaller) installWTAApplication(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(installer.buttonTemplate()))
}

// TODO: use the client_id and scopes from the installer
func (installer *appInstaller) buttonTemplate() string {
	format := `
	<body>
		Add this to Slack <br>
		<a href="https://slack.com/oauth/authorize?client_id=%s&scope=%s">
			<img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcset="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x" />
		</a>
	</body>
	`
	button := fmt.Sprintf(format, installer.slackClientID, defaultScopes)
	return button
}

// code=75950428352.351913834208.ab4422d8cec8e7b134b8dbe7659e097cecb294122cfebee8456cad92f03d1732&state=
func (installer *appInstaller) oAuthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	//use `code` to get `OAuthResponse` response back
	oauthResponse, err := slack.GetOAuthResponseContext(context.Background(), installer.slackClientID, installer.slackClientSecret, code, installer.redirectURL(), false)

	message := "Something wrong happened"
	if err == nil {
		installer.OAuthResponse = *oauthResponse
		dbClient.SaveToDB(installer)

		message = ""
		message += fmt.Sprintf("\nError?: %v", err)
		message += fmt.Sprintf("\nBot token: %v", oauthResponse.Bot)
		message += fmt.Sprintf("\naccessToken: %v", oauthResponse.AccessToken)
		message += fmt.Sprintf("\nscopes: %v", oauthResponse.Scope)
	}
	w.Write([]byte(message))
}

func (installer *appInstaller) saveToDB() {
	// fmt.Println
}

func (installer *appInstaller) redirectURL() string {
	return "https://" + installer.appURL + "/oauthRedirect"
}
