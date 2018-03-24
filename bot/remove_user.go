package bot

import (
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
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:     "username",
				HelpText: "the Slack handle of the user to remove from a team",
				Format:   anyRegex,
			},
			cmd.Argument{
				Name:      "team",
				HelpText:  "the team to remove the user from",
				Format:    anyRegex,
				MultiWord: true,
			},
		},
		HandleFunc: ch,
	}
}

// removeUser removes a user from a team.
func (b *Bot) removeUser(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}

	team := model.Team{
		Name: c.Args[1].Value,
	}
	if err := b.dal.GetTeamByName(&team); err != nil {
		log.WithError(err).Error("Failed to get team " + team.Name)
		return "Failed to get team " + team.Name, noParams
	}

	memberSlackID := parseMention(c.Args[0].Value)
	member := model.Member{
		SlackID: memberSlackID,
	}
	if err := b.dal.GetMemberBySlackID(&member); err != nil {
		log.WithError(err).Error("Failed to get member " + c.Args[0].Value)
		return "Failed to get member " + c.Args[0].Value, noParams
	}

	// Remove user from GitHub team
	if err := b.gh.RemoveUserFromTeam(member.GithubUsername, team.GithubTeamID); err != nil {
		log.WithError(err).Errorf("Failed to remove member %s from GitHub team %s",
			c.Args[0].Value, team.Name)
		return "Failed to remove member from GitHub team", noParams
	}

	teamMember := model.TeamMember{
		MemberSlackID: memberSlackID,
		GithubTeamID:  team.GithubTeamID,
	}
	// Remove user team relation from DB
	if err := b.dal.DeleteTeamMember(&teamMember); err != nil {
		log.WithError(err).Error("Failed to remove member " +
			c.Args[0].Value + " from team " + c.Args[1].Value)
		return "Failed to remove member from team", noParams
	}
	return toMention(member.SlackID) +
		" was removed from `" + team.Name + "` :tada:", noParams
}
