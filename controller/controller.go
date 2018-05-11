package controller

import (
	"SlackPlatform/crossfunction"
	"SlackPlatform/models"

	"github.com/nlopes/slack"
)

var (
	installer      *appInstaller
	dbClient       crossfunction.DBClient
	baristaCommand slashCommand
	interact       interactive
)

//To be passed threw the `StartupControllers`
const (
	appURL            = "goplatform.ngrok.io"
	slackClientID     = "75950428352.351830378721"
	slackClientSecret = "a56df86a6f1fae41f4efceaf20fb9842"
	verificationToken = "8ycguzKPPcWvt7wIsud0a9EL"
)

// StartupControllers call this function to setup the controllers
func StartupControllers(dbWrapper *models.DBWrapper, slackAPI *slack.Client) {
	dbClient = dbWrapper

	installer := &appInstaller{
		appURL:            appURL,
		slackClientID:     slackClientID,
		slackClientSecret: slackClientSecret,
		verificationToken: verificationToken,
	}
	installer.registerRoutes()

	baristaCommand := slashCommand{
		name: "coffeeCommand",
		api:  slackAPI,
	}
	baristaCommand.registerRoutes()

	interact := interactive{}
	interact.registerRoutes()
}
