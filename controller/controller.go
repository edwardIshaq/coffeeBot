package controller

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
	baristaCommand = slashCommand{"coffeeCommand", nil}
	interact = interactive{}

}

// StartupControllers call this function to setup the controllers
func StartupControllers(gormDB *gorm.DB, slackAPI *slack.Client) {
	db = gormDB
	//Demo routes
	registerHelloRoute()
	registerOutgoingHookRoute()

	installer.registerRoutes()
	baristaCommand.registerRoutes()
	interact.registerRoutes()
	userScopesDemo.registerRoutes()
}
