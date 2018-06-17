package core

import (
	"fmt"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewRemoveUserCmd returns a remove user command that removes a user
func NewRemoveUserCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "remove-user",
		HelpText: "Remove a user from a team",
		Options: map[string]*cmd.Option{
			"user": &cmd.Option{
				Key:      "user",
				HelpText: "the Slack handle of the user to remove from a team",
				Format:   cmd.AnyRegex,
			},
			"team": &cmd.Option{
				Key:      "team",
				HelpText: "the team to remove the user from",
				Format:   cmd.AnyRegex,
				Required: true,
			},
		},
		HandleFunc: ch,
	}
}

// removeUser removes a user from a team.
func (core *Plugin) removeUser(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}

	teamName := c.Options["team"].Value
	username := c.Options["user"].Value
	team := model.Team{
		Name: teamName,
	}
	if err := core.Bot.DAL.GetTeamByName(&team); err != nil {
		log.WithError(err).Error("Failed to get team " + team.Name)
		return "Failed to get team " + team.Name, noParams
	}

	memberSlackID := cmd.ParseMention(username)
	member := model.Member{
		SlackID: memberSlackID,
	}
	if err := core.Bot.DAL.GetMemberBySlackID(&member); err != nil {
		log.WithError(err).Error("Failed to get member " + username)
		return "Failed to get member " + username, noParams
	}

	// Remove user from GitHub team
	if err := core.Bot.GitHub.RemoveUserFromTeam(member.GithubUsername, team.GithubTeamID); err != nil {
		log.WithError(err).Errorf("Failed to remove member %s from GitHub team %s",
			member.Name, team.Name)
		msg := fmt.Sprintf("Failed to remove user %s from GitHub team %s. "+
			"Make sure %s's GitHub ID (currently \"%s\") is correct.",
			member.Name, team.Name, member.Name, member.GithubUsername)
		return msg, noParams
	}

	teamMember := model.TeamMember{
		MemberSlackID: memberSlackID,
		GithubTeamID:  team.GithubTeamID,
	}
	// Remove user team relation from DB
	if err := core.Bot.DAL.DeleteTeamMember(&teamMember); err != nil {
		log.WithError(err).Error("Failed to remove member " +
			member.Name + " from team " + team.Name)
		return "Failed to remove member from team", noParams
	}
	return cmd.ToMention(member.SlackID) +
		" was removed from `" + team.Name + "` :tada:", noParams
}
