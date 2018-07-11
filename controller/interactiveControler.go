package controller

import (
	"fmt"
	"net/http"
	"regexp"
)

type interactiveController struct {
	// every interactive component needs a callbackID
	// to route the calls back from the interactive handler
	callbackID string

	//Regex to match against callbacks
	callbackRegex *regexp.Regexp

	//Next interaction
	nextController *interactiveController
}

func simpleInteractiveController(callback string) *interactiveController {
	return newInteractiveController(callback, callback)
}

func newInteractiveController(callback, pattern string) *interactiveController {
	callbackRegex, _ := regexp.Compile(pattern)
	return &interactiveController{
		callbackID:    callback,
		callbackRegex: callbackRegex,
	}
}

func (i *interactiveController) canHandleCallback(callback string) bool {
	return i.callbackRegex.MatchString(callback)
}

func (i *interactiveController) handleCallback(w http.ResponseWriter, r *http.Request, actionCallback SlackActionCallback) {
	if !processSlackClient(w, r) {
		return
	}
	r.ParseForm()

	fmt.Println("interactiveController controller")
	fmt.Println(actionCallback)
}

func processSlackClient(w http.ResponseWriter, r *http.Request) bool {
	_, ok := getSlackClientFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
	}
	return ok
}
