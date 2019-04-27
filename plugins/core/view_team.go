package core

import (
	"fmt"

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
func (core *Plugin) viewTeam(c cmd.Context) (string, slack.PostMessageParameters) {
	params := slack.PostMessageParameters{}
	team := &model.Team{
		Name: c.Options["team"].Value,
	}

	// Fetch team from DB
	if err := core.Bot.DAL.GetTeamByName(team); err != nil {
		log.WithError(err).Errorf("Failed to get team %s", team.Name)
		return "Failed to get team " + team.Name, params
	}

	// Fetch team tech leads
	techLeads, err := core.Bot.DAL.GetTechLeadsByTeam(team)
	if err != nil {
		log.WithError(err).Errorf("Failed to get tech leads for team %s", team.Name)
		return "Failed to get tech leads for team " + team.Name, params
	}

	params.Attachments = team.SlackAttachments(*techLeads)

	// Fetch GitHub team name since we don't store it in the DB
	if ghTeam, err := core.Bot.GitHub.GetTeam(team.GithubTeamID); err == nil {
		ghNameAttachment := slack.Attachment{
			Text:  "GitHub Team Name: " + *ghTeam.Name,
			Color: "good",
		}
		params.Attachments = append(params.Attachments, ghNameAttachment)
	} else {
		msg := fmt.Sprintf("Failed to find GitHub team with ID %d", team.GithubTeamID)
		core.Bot.Log.WithError(err).Error(msg)
		return "Found team " + team.Name +
			", but an error occurred while fetching its associated GitHub team: " +
			msg, params
	}

	return "Team " + team.Name, params
}
