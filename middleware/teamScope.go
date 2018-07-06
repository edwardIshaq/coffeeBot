package middleware

import (
	"SlackPlatform/crossfunction"
	"context"
	"net/http"

	"github.com/edwardIshaq/slack"
)

// TeamScope middleWare will try to fetch a team for the request
type TeamScope struct {
	Next http.Handler
}

type contextKey string

const (
	contextAccessTokenKey         = contextKey("TeamScope.context.accessTokenKey")
	contextSlackAPIKey            = contextKey("TeamScope.context.slackAPI")
	contextStagingChannelIDKey    = contextKey("TeamScope.context.stagingChannelID")
	contextProductionChannelIDKey = contextKey("TeamScope.context.productionChannelID")
)

func (mw *TeamScope) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if mw.Next == nil {
		mw.Next = http.DefaultServeMux
	}

	team := crossfunction.TeamForRequest(r)
	if team != nil && len(team.AccessToken) > 0 {
		ctx := context.WithValue(r.Context(), contextAccessTokenKey, team.AccessToken)
		ctx = context.WithValue(ctx, contextStagingChannelIDKey, team.StagingChannelID)
		ctx = context.WithValue(ctx, contextProductionChannelIDKey, team.ProductionChannelID)
		slackAPI := slack.New(team.AccessToken)
		slackAPI.SetDebug(true)
		ctx = context.WithValue(ctx, contextSlackAPIKey, slackAPI)
		tokenContext := r.WithContext(ctx)
		mw.Next.ServeHTTP(w, tokenContext)
	} else {
		mw.Next.ServeHTTP(w, r)
	}
}

// AccessToken gets the team's AccessToken from the context.
func AccessToken(ctx context.Context) (string, bool) {
	tokenStr, ok := ctx.Value(contextAccessTokenKey).(string)
	return tokenStr, ok
}

// SlackAPI returns the slackAPI client
func SlackAPI(ctx context.Context) (*slack.Client, bool) {
	slackAPI, ok := ctx.Value(contextSlackAPIKey).(*slack.Client)
	return slackAPI, ok
}

// StagingChannelID extracts stagingChannelID from the context
func StagingChannelID(ctx context.Context) (string, bool) {
	stagingChannelID, ok := ctx.Value(contextStagingChannelIDKey).(string)
	ok = len(stagingChannelID) > 0 && ok
	return stagingChannelID, ok
}

// ProductionChannelID extracts stagingChannelID from the context
func ProductionChannelID(ctx context.Context) (string, bool) {
	prodChannelID, ok := ctx.Value(contextProductionChannelIDKey).(string)
	ok = len(prodChannelID) > 0 && ok
	return prodChannelID, ok
}
