package bot

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
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:      "team-name",
				HelpText:  "the name of the team to view",
				Format:    anyRegex,
				MultiWord: true,
			},
		},
		HandleFunc: ch,
	}
}

// viewTeam displays a teams's information.
func (b *Bot) viewTeam(c cmd.Context) (string, slack.PostMessageParameters) {
	params := slack.PostMessageParameters{}
	team := model.Team{
		Name: c.Args[0].Value,
	}
	if err := b.dal.GetTeamByName(&team); err != nil {
		log.WithError(err).Error("Failed to get team " + team.Name)
		return "Failed to get team " + team.Name, params
	}
	params.Attachments = team.SlackAttachments()
	return "Team " + c.Args[0].Value, params
}