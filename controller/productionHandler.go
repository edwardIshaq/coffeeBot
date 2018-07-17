package controller

import (
	"SlackPlatform/models"
	"net/http"
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
}
