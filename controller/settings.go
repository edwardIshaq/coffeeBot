package controller

import (
	"SlackPlatform/crossfunction"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/edwardIshaq/slack"
)

const (
	keyStagingChannel    = "staging_channel"
	keyProductionChannel = "production_channel"
)

type settings struct {
	// messageMenu will generate an interactive component with `callbackID` callback
	callbackID string
	//Regex to match against callbacks
	callbackRegex *regexp.Regexp
}

func newSettings() *settings {
	callbackID := "settings_callbackID"
	regex, _ := regexp.Compile(callbackID)
	return &settings{
		callbackID:    callbackID,
		callbackRegex: regex,
	}
}
func (s *settings) sendSettingsDialog(channelID, trigger string) {
	channelsDialog := slack.Dialog{
		TriggerID:   trigger,
		CallbackID:  s.callbackID,
		Title:       "CoffeeBot Settings",
		SubmitLabel: "Configure",
	}

	stagingInput := slack.NewChannelsSelect(keyStagingChannel, "Staging Channel")
	productionInput := slack.NewChannelsSelect(keyProductionChannel, "Production Channel")
	channelsDialog.Elements = []slack.DialogElement{stagingInput, productionInput}
	api.OpenDialog(channelsDialog)
}

func (s *settings) canHandleCallback(callback string) bool {
	return s.callbackRegex.MatchString(callback)
}

func (s *settings) handleCallback(w http.ResponseWriter, r *http.Request, actionCallback SlackActionCallback) {
	_, ok := getSlackClientFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dialogResponse := slack.DialogSubmitCallback{}
	json.Unmarshal([]byte(r.PostFormValue("payload")), &dialogResponse)
	if len(dialogResponse.Submission) > 0 {
		stagingChannelID := dialogResponse.Submission[keyStagingChannel]
		productionChannelID := dialogResponse.Submission[keyProductionChannel]

		if len(stagingChannelID) > 0 && len(productionChannelID) > 0 {
			team := crossfunction.TeamForRequest(r)
			team.UpdateChannels(stagingChannelID, productionChannelID)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}
