package models

import (
	"github.com/jinzhu/gorm"
	"github.com/nlopes/slack"
)

// Team is the slack representation of Workspace
type Team struct {
	gorm.Model

	TeamID      string
	TeamName    string
	AccessToken string
	Scope       string
	UserID      string

	BotUserID      string
	BotAccessToken string
}

// NewTeam constructs a new `Team` with `slack.OAuthResponse`
func NewTeam(oauth *slack.OAuthResponse) *Team {
	team := &Team{
		TeamID:         oauth.TeamID,
		TeamName:       oauth.TeamName,
		AccessToken:    oauth.AccessToken,
		Scope:          oauth.Scope,
		UserID:         oauth.UserID,
		BotUserID:      oauth.Bot.BotUserID,
		BotAccessToken: oauth.Bot.BotAccessToken,
	}
	//TODO: insert/update existing team
	db.Debug().Save(team)
	return team
}

// TeamByID returns a team if found or Nil
func TeamByID(teamID string) *Team {
	team := &Team{}
	db.Debug().Where("team_id = ?", teamID).First(&team)
	return team
}
