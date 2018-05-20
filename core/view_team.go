package core

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewViewTeamCmd returns a view team command that displays information about a user
func NewViewTeamCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "view-team",
		HelpText: "View information about a Launch Pad team",
		Options: map[string]*cmd.Option{
			"team": &cmd.Option{
				Key:      "team",
				HelpText: "the name of the team to view",
				Format:   cmd.AnyRegex,
				Required: true,
			},
		},
		HandleFunc: ch,
	}
}

// viewTeam displays a teams's information.
func (core *CorePlugin) viewTeam(c cmd.Context) (string, slack.PostMessageParameters) {
	params := slack.PostMessageParameters{}
	team := model.Team{
		Name: c.Options["team"].Value,
	}
	if err := core.Bot.DAL.GetTeamByName(&team); err != nil {
		log.WithError(err).Error("Failed to get team " + team.Name)
		return "Failed to get team " + team.Name, params
	}
	params.Attachments = team.SlackAttachments()
	return "Team " + team.Name, params
}
