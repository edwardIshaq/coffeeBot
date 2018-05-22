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
	newTeam := &Team{}
	db.Where(&Team{TeamID: oauth.TeamID}).FirstOrInit(&newTeam)

	newTeam.TeamName = oauth.TeamName
	newTeam.AccessToken = oauth.AccessToken
	newTeam.Scope = oauth.Scope
	newTeam.UserID = oauth.UserID
	newTeam.BotUserID = oauth.Bot.BotUserID
	newTeam.BotAccessToken = oauth.Bot.BotAccessToken
	db.Save(&newTeam)
	return newTeam
}

// TeamByID returns a team if found or Nil
func TeamByID(teamID string) *Team {
	team := &Team{}
	db.Debug().Where("team_id = ?", teamID).First(&team)
	return team
}
