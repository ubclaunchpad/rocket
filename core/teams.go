package core

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewTeamsCmd returns a teams command that displays a list of Launch Pad teams
func NewTeamsCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:       "teams",
		HelpText:   "List Launch Pad teams",
		Options:    map[string]*cmd.Option{},
		HandleFunc: ch,
	}
}

// listTeams displays Launch Pad teams
func (core *CorePlugin) listTeams(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	teams := model.Teams{}
	if err := core.Bot.DAL.GetTeamNames(&teams); err != nil {
		log.WithError(err).Error("Failed to get team names")
		return "Failed to get team names", noParams
	}
	names := ""
	for _, team := range teams {
		names += team.Name + "\n"
	}
	return names, noParams
}
