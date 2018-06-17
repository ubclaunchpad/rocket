package core

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewEditTeamCmd returns an add team command that creates a new Launch Pad team
func NewEditTeamCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "edit-team",
		HelpText: "Update an existing Launch Pad team (admins and tech leads only)",
		Options: map[string]*cmd.Option{
			"team": &cmd.Option{
				Key:      "team",
				HelpText: "the name of the existing team",
				Format:   cmd.AnyRegex,
				Required: true,
			},
			"name": &cmd.Option{
				Key:      "name",
				HelpText: "the new name of the team",
				Format:   cmd.AnyRegex,
				Required: false,
			},
			"platform": &cmd.Option{
				Key:      "platform",
				HelpText: "the new platform the team develops on (i.e iOS, Android etc)",
				Format:   cmd.AnyRegex,
				Required: false,
			},
		},
		HandleFunc: ch,
	}
}

// editTeam edits an existing Launch Pad team.
func (core *Plugin) editTeam(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	if !c.User.IsAdmin && !c.User.IsTechLead {
		return "You must be an admin or tech lead to use this command", noParams
	}

	currentName := c.Options["team"].Value
	currentTeam := &model.Team{
		Name: currentName,
	}
	newTeam := &model.Team{
		Name:     c.Options["name"].Value,
		Platform: c.Options["platform"].Value,
	}

	// Try get existing team from DB
	if err := core.Bot.DAL.GetTeamByName(currentTeam); err != nil {
		core.Bot.Log.WithError(err).Errorf("failed to update team %s", currentName)
		return "Failed to update team " + currentName, noParams
	}

	// Finally, update team in DB
	if err := core.Bot.DAL.UpdateTeam(currentTeam, newTeam); err != nil {
		log.WithError(err).Errorf("failed to update team %s", currentName)
		return "Failed to update team " + currentName, noParams
	}

	return "`" + currentName + "` has been updated :tada:", noParams
}
