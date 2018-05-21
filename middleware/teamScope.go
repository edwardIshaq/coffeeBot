package middleware

import (
	"SlackPlatform/crossfunction"
	"context"
	"net/http"
)

// TeamScope middleWare will try to fetch a team for the request
type TeamScope struct {
	Next http.Handler
}

type contextKey string

const (
	contextAccessTokenKey = contextKey("TeamScope.context.accessTokenKey")
)

func (mw *TeamScope) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if mw.Next == nil {
		mw.Next = http.DefaultServeMux
	}

	team := crossfunction.TeamForRequest(r)
	if len(team.AccessToken) > 0 {
		ctx := context.WithValue(r.Context(), contextAccessTokenKey, team.AccessToken)
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
