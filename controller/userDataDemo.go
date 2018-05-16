package controller

import (
	"SlackPlatform/crossfunction"
	"SlackPlatform/models"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/nlopes/slack"
)

const (
	userDataScopes        = "reminders:write:user"
	permissionUserDataAPI = "apps.permissions.users.request"
)

type userDataDemo struct {
	scopes string
}

func newUserDataDemo() *userDataDemo {
	return &userDataDemo{
		scopes: userDataScopes,
	}
}

func (userData *userDataDemo) registerRoutes() {
	http.HandleFunc("/request_userdata_scopes", userData.handleUserDataRequest)
}

func (userData *userDataDemo) handleUserDataRequest(w http.ResponseWriter, r *http.Request) {
	triggerID := r.PostFormValue("trigger_id")
	userID := r.PostFormValue("user_id")
	teamID := r.PostFormValue("team_id")
	team := models.TeamByID(teamID)
	slackClient := crossfunction.ClientForRequest(r)

	userData.sendAPIRequest(userID, triggerID, team, slackClient)
}

func (userData *userDataDemo) sendAPIRequest(userID, triggerID string, team *models.Team, client *slack.Client) {
	postBody := url.Values{
		"token":      {team.AccessToken},
		"trigger_id": {triggerID},
		"scopes":     {userData.scopes},
	}
	reqBody := strings.NewReader(postBody.Encode())
	req, err := http.NewRequest("POST", slack.SLACK_API+permissionUserDataAPI, reqBody)
	if err != nil {
		fmt.Println("error happened: ", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(context.Background())

	//Fire the request
	resp, err := slack.HTTPClient.Do(req)
	if err != nil {
		fmt.Println("\nResponseError: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Printf("error reading body %v", err2)
			return
		}
		bodyString := string(bodyBytes)
		fmt.Println("\nbodyString: ", bodyString)
	}
}
