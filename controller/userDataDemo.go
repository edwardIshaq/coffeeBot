package controller

import (
	"SlackPlatform/crossfunction"
	"SlackPlatform/middleware"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/nlopes/slack"
)

const (
	channelScopes         = "channels:history"
	teamScopes            = "emoji:read"
	permissionAddScopeAPI = "apps.permissions.request"

	userScopes            = "reminders:write:user"
	permissionUserDataAPI = "apps.permissions.users.request"
)

func registerPermissionsRequestsRoutes() {
	//invoked via slash command: `/request_user_data`
	http.HandleFunc("/request_userdata_scopes", handleScopeRequests)
	//invoked via slash command: `/request_team_data`
	http.HandleFunc("/request_team_scopes", handleScopeRequests)
	//invoked via slash command: `/request_channel_data`
	http.HandleFunc("/request_channel_scopes", handleScopeRequests)
}

func handleScopeRequests(w http.ResponseWriter, r *http.Request) {
	fmt.Println()
	fmt.Printf("called handleScopeRequests %s\n", r.RequestURI)
	fmt.Printf("Passed PostForm: %v", r.PostForm)

	var scopes string
	var apiURL string

	switch r.RequestURI {
	case "/request_team_scopes":
		scopes = teamScopes
		apiURL = permissionAddScopeAPI

	case "/request_channel_scopes":
		scopes = channelScopes
		apiURL = permissionAddScopeAPI

	case "/request_userdata_scopes":
		scopes = userScopes
		apiURL = permissionUserDataAPI
	}

	fmt.Println()
	fmt.Printf("request scopes: %s, using api: %s", scopes, apiURL)
	fmt.Println()

	triggerID := r.PostFormValue("trigger_id")
	userID := r.PostFormValue("user_id")

	api = crossfunction.ClientForRequest(r)
	accessToken, ok := middleware.AccessToken(r.Context())
	if !ok {
		fmt.Printf("Error retrieving access token")
		return
	}
	requestScopes(userID, triggerID, accessToken, scopes, apiURL)
}

func requestScopes(userID, triggerID, accessToken, scopes, endpoint string) {
	postBody := url.Values{
		"token":       {accessToken},
		"scopes":      {scopes},
		"trigger_id":  {triggerID},
		"did_confirm": {"false"},
	}
	fmt.Printf("\npostBody = %v\n", postBody)
	reqBody := strings.NewReader(postBody.Encode())
	req, err := http.NewRequest("POST", slack.SLACK_API+endpoint, reqBody)
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

	// if resp.StatusCode == http.StatusOK
	{
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Printf("error reading body %v", err2)
			return
		}
		bodyString := string(bodyBytes)
		fmt.Println("\nbodyString: ", bodyString)
	}
}
