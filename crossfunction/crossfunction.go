package crossfunction

import (
	"SlackPlatform/models"
	"net/http"

	"github.com/nlopes/slack"
)

// ClientForRequest finds the teamID from request and tries to return a client
func ClientForRequest(r *http.Request) *slack.Client {
	teamID := r.PostFormValue("team_id")
	team := models.TeamByID(teamID)
	api := slack.New(team.AccessToken)
	api.SetDebug(true)
	return api
}
