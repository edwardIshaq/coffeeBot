package controller

/*
This will support a basic play grounds for the `utaApp` found at the link below.string

add to slack button:	/addToSlack
OAuth handlers:			/oauthRedirect
outgoing hook url: 		/outgoingHooks
/barista command url:	/coffeeCommand
interactive url:		/interactive
Dynamic Menu:			/dynamicMenu

TODO:
------------------------------------------
[√] save tokens to DB
[ ] Add runtime flag to switch between dev/prod
*/

import (
	"log"

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

	log.Printf("Before installer register routes: key-tail= %s", GetAppSecretTail())
	installer.registerRoutes()
	baristaCommand.registerRoutes()
	interact.registerRoutes()
	userScopesDemo.registerRoutes()
}
