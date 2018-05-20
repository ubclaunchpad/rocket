package core

import (
	"fmt"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewAddUserCmd returns an add command that adds a user
func NewAddUserCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "add-user",
		HelpText: "Add a user to a team",
		Options: map[string]*cmd.Option{
			"user": &cmd.Option{
				Key:      "user",
				HelpText: "the Slack handle of the user to add to a team",
				Format:   cmd.AnyRegex,
				Required: true,
			},
			"team": &cmd.Option{
				Key:      "team",
				HelpText: "the team to add the user to",
				Format:   cmd.AnyRegex,
				Required: true,
			},
		},
		HandleFunc: ch,
	}
}

// addUser adds an existing user to a team.
func (core *CorePlugin) addUser(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	username := c.Options["user"].Value
	teamName := c.Options["team"].Value

	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}

	team := model.Team{
		Name: teamName,
	}
	if err := core.Bot.DAL.GetTeamByName(&team); err != nil {
		log.WithError(err).Error("Failed to find team " + team.Name)
		return "Failed to find team " + team.Name, noParams
	}

	slackID := cmd.ParseMention(username)
	member := model.Member{
		SlackID: slackID,
	}
	if err := core.Bot.DAL.GetMemberBySlackID(&member); err != nil {
		log.WithError(err).Errorf("Failed to find member %s", username)
		return "Failed to find member " + username, noParams
	}

	// Add user to corresponding GitHub team
	if err := core.Bot.GitHub.AddUserToTeam(member.GithubUsername, team.GithubTeamID); err != nil {
		log.WithError(err).Errorf("Failed to add user %s to GitHub team %s",
			member.Name, team.Name)
		msg := fmt.Sprintf("Failed to add user %s to GitHub team %s. "+
			"Make sure %s's GitHub ID (currently \"%s\") is correct.",
			member.Name, team.Name, member.Name, member.GithubUsername)
		return msg, noParams
	}

	teamMember := model.TeamMember{
		MemberSlackID: slackID,
		GithubTeamID:  team.GithubTeamID,
	}
	// Finally, add relation to DB
	if err := core.Bot.DAL.CreateTeamMember(&teamMember); err != nil {
		log.WithError(err).Errorf("Failed to add member %s to team %s",
			member.Name, team.Name)
		return fmt.Sprintf("Failed to add member %s to team %s",
			member.Name, team.Name), noParams
	}
	return cmd.ToMention(member.SlackID) +
		" was added to `" + team.Name + "` team :tada:", noParams
}
