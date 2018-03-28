package bot

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewRemoveTeamCmd returns a remove team command that removes a new Launch Pad team
func NewRemoveTeamCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "remove-team",
		HelpText: "Delete a new Launch Pad team",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:      "team-name",
				HelpText:  "the name of the team to remove",
				Format:    anyRegex,
				MultiWord: true,
			},
		},
		HandleFunc: ch,
	}
}

// removeTeam removes a Launch Pad team.
func (b *Bot) removeTeam(c cmd.Context) (string, slack.PostMessageParameters) {
	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}

	noParams := slack.PostMessageParameters{}
	team := model.Team{
		Name: c.Args[0].Value,
	}
	if err := b.dal.GetTeamByName(&team); err != nil {
		log.WithError(err).Error("Failed to find team " + team.Name)
		return "Failed to find team " + team.Name, noParams
	}

	// Remove team from GitHub
	if err := b.gh.RemoveTeam(team.GithubTeamID); err != nil {
		log.WithError(err).Error("Failed to remove GitHub team " + team.Name)
		return "Failed to remove GitHub team " + team.Name, noParams
	}

	// Finally remove team from database
	if err := b.dal.DeleteTeamByName(&team); err != nil {
		log.WithError(err).Error("Failed to delete team " + team.Name)
		return "Failed to delete team " + team.Name, noParams
	}
	return "`" + team.Name + "` team has been deleted :tada:", noParams
}
