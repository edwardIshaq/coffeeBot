package crossfunction

import (
	"SlackPlatform/models"
	"net/http"

	"github.com/nlopes/slack"
)

// TeamForRequest finds the matching team for the `team_id`` post param
func TeamForRequest(r *http.Request) *models.Team {
	teamID := r.PostFormValue("team_id")
	team := models.TeamByID(teamID)
	return team
}

// ClientForRequest finds the teamID from request and tries to return a client
func ClientForRequest(r *http.Request) *slack.Client {
	team := TeamForRequest(r)
	api := slack.New(team.AccessToken)
	api.SetDebug(true)
	return api
}
