package crossfunction

import (
	"SlackPlatform/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nlopes/slack"
)

// teamIDFromRequest Extracts teamID from a slack request
// and fetches the team from db
func teamIDFromRequest(r *http.Request) (string, error) {
	teamID := r.PostFormValue("team_id")
	if len(teamID) > 0 {
		return teamID, nil
	}

	//else decode the body to find the team_id
	slackResp := &models.PayloadResponse{}
	payload := r.PostFormValue("payload")
	json.Unmarshal([]byte(payload), &slackResp)

	responseTeamID := slackResp.Team.ID
	if len(responseTeamID) > 0 {
		return responseTeamID, nil
	}
	return "", fmt.Errorf("couldnt find any team_id")
}

// TeamForRequest finds the matching team for the `team_id`` post param
func TeamForRequest(r *http.Request) *models.Team {
	teamID, err := teamIDFromRequest(r)
	if err != nil {
		return nil
	}
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
