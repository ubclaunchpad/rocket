package core

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
		HelpText: "Create a new Launch Pad team (admins only)",
		Options: map[string]*cmd.Option{
			"name": &cmd.Option{
				Key:      "name",
				HelpText: "the name of the new team",
				Format:   cmd.AnyRegex,
				Required: true,
			},
			"platform": &cmd.Option{
				Key:      "platform",
				HelpText: "the platform the team develops on (i.e iOS, Android etc)",
				Format:   cmd.AnyRegex,
				Required: true,
			},
			"github": &cmd.Option{
				Key:      "github",
				HelpText: "the name of the team to create on GitHub",
				Format:   cmd.AnyRegex,
				Required: false,
			},
		},
		HandleFunc: ch,
	}
}

// addTeam creates a new Launch Pad team.
func (core *Plugin) addTeam(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}

	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}

	teamName := c.Options["name"].Value
	platform := c.Options["platform"].Value

	// Set custom GitHub team name if applicable
	ghTeamName := strings.ToLower(strings.Replace(teamName, " ", "-", -1))
	if c.Options["github"].Value != "" {
		ghTeamName = c.Options["github"].Value
	}

	// Create the team on GitHub
	ghTeam, err := core.Bot.GitHub.CreateTeam(ghTeamName)
	core.Bot.Log.Info("create team, ", ghTeam, err)
	if err != nil {
		log.WithError(err).Errorf("Failed to create team %s on GitHub", teamName)
		return "Failed to create team " + teamName + " on GitHub", noParams
	}

	team := model.Team{
		Name:         teamName,
		Platform:     platform,
		GithubTeamID: int(*ghTeam.ID),
	}
	// Finally, add team to DB
	if err := core.Bot.DAL.CreateTeam(&team); err != nil {
		log.WithError(err).Errorf("Failed to create team %s", team.Name)
		return "Failed to create team " + team.Name, noParams
	}

	return "`" + team.Name + "` has been added :tada:", noParams
}
