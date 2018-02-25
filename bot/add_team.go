package bot

import (
	"strings"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewAddTeamCmd returns an add team command that creates a new Launch Pad team
func NewAddTeamCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "add-team",
		HelpText: "Create a new Launch Pad team",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:      "team-name",
				HelpText:  "the name of the new team",
				Format:    anyRegex,
				MultiWord: true,
			},
		},
		HandleFunc: ch,
	}
}

// addTeam creates a new Launch Pad team.
func (b *Bot) addTeam(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}

	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}

	teamName := c.Args[0].Value
	// teamName = "Great Team", ghTeamName = "great-team"
	ghTeamName := strings.ToLower(strings.Replace(teamName, " ", "-", -1))

	// Create the team on GitHub
	ghTeam, err := b.gh.CreateTeam(ghTeamName)
	b.log.Info("create team, ", ghTeam, err)
	if err != nil {
		log.WithError(err).Errorf("Failed to create team %s on GitHub", teamName)
		return "Failed to create team " + teamName + " on GitHub", noParams
	}

	team := model.Team{
		Name:         teamName,
		GithubTeamID: *ghTeam.ID,
	}
	// Finally, add team to DB
	if err := b.dal.CreateTeam(&team); err != nil {
		log.WithError(err).Errorf("Failed to create team %s", team.Name)
		return "Failed to create team " + team.Name, noParams
	}
	return "`" + team.Name + "` has been added :tada:", noParams
}
