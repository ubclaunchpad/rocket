package core

import (
	"fmt"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewToggleTechLeadCmd returns a toggle tech lead command that makes an
// existing user the teach lead of a given team if they are not already
// a tech lead of that team, otherwise it removes them from the tech leads role
// on that team.
func NewToggleTechLeadCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "toggle-tech-lead",
		HelpText: "Make an existing user a tech lead of an existing team (admins only)",
		Options: map[string]*cmd.Option{
			"user": &cmd.Option{
				Key:      "user",
				HelpText: "the Slack handle of the user to update",
				Format:   cmd.AnyRegex,
				Required: true,
			},
			"team": &cmd.Option{
				Key:      "team",
				HelpText: "the name of the team on which to set the new tech lead",
				Format:   cmd.AnyRegex,
				Required: true,
			},
		},
		HandleFunc: ch,
	}
}

// toggleTechLead makes an existing user the teach lead of a given team if they
// are not already a tech lead of that team, otherwise it removes them from
// the tech leads role on that team.
func (core *Plugin) toggleTechLead(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}

	// Fetch the member
	username := c.Options["user"].Value
	member := &model.Member{SlackID: cmd.ParseMention(username)}
	if err := core.Bot.DAL.GetMemberBySlackID(member); err != nil {
		log.WithError(err).Errorf("Failed to get %s", username)
		return "Failed to find user", noParams
	}

	// Fetch team
	teamName := c.Options["team"].Value
	team := &model.Team{Name: teamName}
	if err := core.Bot.DAL.GetTeamByName(team); err != nil {
		log.WithError(err).Errorf("Failed to get %s", teamName)
		return "Failed to find team", noParams
	}

	userMention := cmd.ToMention(member.SlackID)

	// Make the member a tech lead of the given team
	teamMember := &model.TeamMember{
		GithubTeamID:  team.GithubTeamID,
		MemberSlackID: member.SlackID,
	}

	// Fetch team member so we can see its updated tech lead status
	if err := core.Bot.DAL.GetTeamMember(teamMember); err != nil {
		log.WithError(err).Errorf("Failed to get team member")
		return "Failed to get team member", noParams
	}

	// Update tech lead status
	teamMember.IsTechLead = !teamMember.IsTechLead
	if err := core.Bot.DAL.SetIsTechLead(teamMember); err != nil {
		log.WithError(err).Errorf("Failed to set %s as tech lead of %s",
			username, teamName)
		return fmt.Sprintf("Failed to set %s as a tech lead of %s. "+
			"Make sure %s if part of %s",
			userMention, teamName, userMention, teamName,
		), noParams
	}

	var status string
	if teamMember.IsTechLead {
		status = "now"
	} else {
		status = "no longer"
	}
	return fmt.Sprintf(
		"%s is %s tech lead of %s :tada:", userMention, status, teamName,
	), noParams
}
