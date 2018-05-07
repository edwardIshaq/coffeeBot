package controller

import (
	"SlackPlatform/models"
)

var (
	installer *appInstaller
	// dbWrapper *models.DBWrapper
	dbClient DBClient
)

//To be passed threw the `StartupControllers`
const (
	appURL            = "goplatform.ngrok.io"
	slackClientID     = "75950428352.351830378721"
	slackClientSecret = "a56df86a6f1fae41f4efceaf20fb9842"
	verificationToken = "8ycguzKPPcWvt7wIsud0a9EL"
)

// StartupControllers call this function to setup the controllers
func StartupControllers(dbWrapper *models.DBWrapper) {
	installer := &appInstaller{
		appURL:            appURL,
		slackClientID:     slackClientID,
		slackClientSecret: slackClientSecret,
		verificationToken: verificationToken,
	}

	dbClient = dbWrapper
	installer.registerRoutes()
}

// DBClient to work with DBWrapper
type DBClient interface {
	SaveToDB(interface{})
}
