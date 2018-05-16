package controller

/*
This will support a basic play grounds for the `utaApp` found at the link below.string

add to slack button:	https://goplatform.ngrok.io/addToSlack
OAuth handlers:			https://goplatform.ngrok.io/oauthRedirect
outgoing hook url: 		https://goplatform.ngrok.io/outgoingHooks
/barista command url:	https://goplatform.ngrok.io/coffeeCommand
interactive url:		https://goplatform.ngrok.io/interactive
Dynamic Menu:			https://goplatform.ngrok.io/dynamicMenu

TODO:
------------------------------------------
[√] save tokens to DB
[ ] Add runtime flag to switch between dev/prod
*/

import (
	"github.com/jinzhu/gorm"
	"github.com/nlopes/slack"
)

var (
	installer      *appInstaller
	baristaCommand slashCommand
	userScopesDemo *userDataDemo
	interact       interactive
	db             *gorm.DB
	api            *slack.Client
)

//To be passed threw the `StartupControllers`
const (
	appURL            = "goplatform.ngrok.io"
	slackClientID     = "75950428352.351830378721"
	slackClientSecret = "a56df86a6f1fae41f4efceaf20fb9842"
	verificationToken = "8ycguzKPPcWvt7wIsud0a9EL"
)

func init() {
	installer = defaultApp()
	userScopesDemo = newUserDataDemo()
	baristaCommand = slashCommand{"coffeeCommand"}
	interact = interactive{}
}

// StartupControllers call this function to setup the controllers
func StartupControllers(gormDB *gorm.DB) {
	db = gormDB
	//Demo routes
	registerHelloRoute()
	registerOutgoingHookRoute()

	installer.registerRoutes()
	baristaCommand.registerRoutes()
	interact.registerRoutes()
	userScopesDemo.registerRoutes()
}
