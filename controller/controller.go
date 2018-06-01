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
	"github.com/jinzhu/gorm"
	"github.com/nlopes/slack"
)

var (
	installer     *appInstaller
	slashBarista  *slashCommand
	dialogHandler *dialogInteraction
	saveBevAction *saveBeverageInteraction
	interact      *interactive
	db            *gorm.DB
	api           *slack.Client
)

func init() {
	installer = defaultApp()
	interact = &interactive{}

	slashBarista = baristaCommand()
	interact.addComponent(slashBarista)

	dialogHandler = beverageDialogInteraction()
	interact.addComponent(dialogHandler)

	saveBevAction = saveBevInteraction()
	interact.addComponent(saveBevAction)
}

// StartupControllers call this function to setup the controllers
func StartupControllers(gormDB *gorm.DB) {
	db = gormDB
	//Demo routes
	registerHelloRoute()
	registerOutgoingHookRoute()
	registerPermissionsRequestsRoutes()

	installer.registerRoutes()
	slashBarista.registerRoutes()
	interact.registerRoutes()
}
