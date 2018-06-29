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
[√] Add runtime flag to switch between dev/prod
*/

import (
	"SlackPlatform/middleware"
	"net/http"

	"github.com/edwardIshaq/slack"
	"github.com/jinzhu/gorm"
)

var (
	installer           *appInstaller
	slashBarista        *slashCommand
	beverageMenuHandler *bevMenuHandler
	dialogHandler       *dialogInteraction
	orderConfirmHandler *orderConfirmInteraction
	saveBevAction       *saveBeverageInteraction
	interact            *interactive
	db                  *gorm.DB
	api                 *slack.Client
)

func init() {
	installer = defaultApp()
	interact = &interactive{}

	slashBarista = baristaCommand()
	// interact.addComponent(slashBarista)

	beverageMenuHandler = newBevMenuHndler()
	interact.addComponent(beverageMenuHandler)

	dialogHandler = beverageDialogInteraction()
	interact.addComponent(dialogHandler)

	saveBevAction = saveBevInteraction()
	interact.addComponent(saveBevAction)

	orderConfirmHandler = newOrderConfirm()
	interact.addComponent(orderConfirmHandler)
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

// getSlackClientFromRequest gets one from the r.Context()
func getSlackClientFromRequest(r *http.Request) (*slack.Client, bool) {
	slackClient, err := middleware.SlackAPI(r.Context())
	api = slackClient
	return slackClient, err
}
