package controller

import (
	"SlackPlatform/models"
	"context"
	"fmt"
	"log"
	"net/http"

	// "github.com/jinzhu/gorm"
	"github.com/nlopes/slack"
)

/*
Handler for the install button
*/
func (installer *appInstaller) registerRoutes() {
	http.HandleFunc("/addToSlack", installer.installWTAApplication)
	http.HandleFunc("/oauthRedirect", installer.oAuthRedirectHandler)
}

// InstallWTAApplication registers a route to install the app on workspaces
func (installer *appInstaller) installWTAApplication(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(installer.buttonTemplate()))
}

// TODO: use the client_id and scopes from the installer
func (installer *appInstaller) buttonTemplate() string {
	format := `
	<body>
		Add this to Slack <br>
		<a href="https://%s/oauth/authorize?client_id=%s&scope=%s">
			<img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcset="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x" />
		</a>
	</body>
	`
	button := fmt.Sprintf(format, installer.slackHost, installer.slackClientID, installer.scopes)
	return button
}

func (installer *appInstaller) oAuthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	//use `code` to get `OAuthResponse` response back
	code := r.FormValue("code")
	oauthResponse, err := slack.GetOAuthResponseContext(context.Background(), installer.slackClientID, installer.slackClientSecret, code, installer.redirectURL(), false)

	message := "Something wrong happened"
	if err == nil {
		installer.OAuthResponse = *oauthResponse
		models.NewTeam(oauthResponse)

		message = ""
		if err != nil {
			message += fmt.Sprintf("\nError?: %v", err)
		}
		message += fmt.Sprintf("\nBot token: %v", oauthResponse.Bot)
		message += fmt.Sprintf("\naccessToken: %v", oauthResponse.AccessToken)
		message += fmt.Sprintf("\nscopes: %v", oauthResponse.Scope)
	}
	log.Printf("OAUTH error: %v", err)
	w.Write([]byte(message))
}
