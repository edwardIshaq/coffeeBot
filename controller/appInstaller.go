package controller

import (
	"github.com/nlopes/slack"
)

func defaultApp() *appInstaller {
	installer = prodApp()
	installer.updateSlackAPI()
	return installer
}

type appInstaller struct {
	slack.OAuthResponse
	slackHost         string
	appURL            string
	slackClientID     string
	slackClientSecret string
	verificationToken string
	scopes            string
}

func (installer *appInstaller) updateSlackAPI() {
	url := "https://" + installer.slackHost + "/api/"
	slack.SLACK_API = url
}

func (installer *appInstaller) redirectURL() string {
	return "https://" + installer.appURL + "/oauthRedirect"
}

func devApp() *appInstaller {
	//Admin page
	//https://api.dev.slack.com/apps/A0QAP1SJD/general
	return &appInstaller{
		appURL:            "goplatform.ngrok.io",
		slackHost:         "dev.slack.com",
		slackClientID:     "8092351665.24363060625",
		slackClientSecret: "f711a81815faa802051475eea0c3874a",
		verificationToken: "mlbDZzxaOiEIZ6I5PIKAwR37",
		scopes:            "commands,chat:write",
	}
}

/*
utaApp info https://api.slack.com/apps/AABQEB4M7
info which is related to a single team `eddie-beta`
Needs to be stored in DB
*/
func prodApp() *appInstaller {
	//Admin page
	// utaApp info https://api.slack.com/apps/AABQEB4M7
	return &appInstaller{
		appURL:            "goplatform.ngrok.io",
		slackHost:         "slack.com",
		slackClientID:     "75950428352.351830378721",
		slackClientSecret: "a56df86a6f1fae41f4efceaf20fb9842",
		verificationToken: "8ycguzKPPcWvt7wIsud0a9EL",
		scopes:            "commands,groups:read,bot,chat:write:bot",
	}
}
