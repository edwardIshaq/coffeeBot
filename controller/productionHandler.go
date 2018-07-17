package controller

import (
	"SlackPlatform/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/edwardIshaq/slack"
)

type productionController struct {
	interactiveController
}

func newProductionHandler() *productionController {
	interactive := simpleInteractiveController("production_handler")
	return &productionController{
		interactiveController: *interactive,
	}
}

func (p *productionController) handleCallback(w http.ResponseWriter, r *http.Request, actionCallback SlackActionCallback) {
	p.interactiveController.handleCallback(w, r, actionCallback)

	orderQuery := models.Order{ProdMsgID: actionCallback.MessageTs}
	fetchedOrder := orderQuery.Fetch()
	fetchedOrder.IsFulfilled = true
	fetchedOrder.Save()

	bevID := strconv.FormatUint(uint64(fetchedOrder.BeverageID), 10)
	beverage := models.BeverageByID(bevID)
	userIDToNotify := beverage.UserID
	postParams := slack.NewPostMessageParameters()
	postParams.Channel = userIDToNotify
	messageText := fmt.Sprintf("your order of `%s` is ready to pickup", beverage.Name)
	api.PostMessage(userIDToNotify, messageText, postParams)
}
