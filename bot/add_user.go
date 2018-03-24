package bot

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
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:      "member",
				HelpText:  "the Slack handle of the user to add to a team",
				Format:    anyRegex,
				MultiWord: false,
			},
			cmd.Argument{
				Name:      "team-name",
				HelpText:  "the team to add the user to",
				Format:    anyRegex,
				MultiWord: true,
			},
		},
		HandleFunc: ch,
	}
}

// addUser adds an existing user to a team.
func (b *Bot) addUser(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}

	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}

	team := model.Team{
		Name: c.Args[1].Value,
	}
	if err := b.dal.GetTeamByName(&team); err != nil {
		log.WithError(err).Error("Failed to find team " + team.Name)
		return "Failed to find team " + team.Name, noParams
	}

	slackID := parseMention(c.Args[0].Value)
	member := model.Member{
		SlackID: slackID,
	}
	if err := b.dal.GetMemberBySlackID(&member); err != nil {
		log.WithError(err).Errorf("Failed to find member %s", c.Args[0].Value)
		return "Failed to find member " + c.Args[0].Value, noParams
	}

	// Add user to corresponding GitHub team
	if err := b.gh.AddUserToTeam(member.GithubUsername, int64(team.GithubTeamID)); err != nil {
		log.WithError(err).Errorf("Failed to add user %s to GitHub team %s",
			c.Args[0].Value, team.Name)
		return fmt.Sprintf("Failed to add user %s to GitHub team %s",
			c.Args[0].Value, team.Name), noParams
	}

	teamMember := model.TeamMember{
		MemberSlackID: slackID,
		GithubTeamID:  team.GithubTeamID,
	}
	// Finally, add relation to DB
	if err := b.dal.CreateTeamMember(&teamMember); err != nil {
		log.WithError(err).Errorf("Failed to add member %s to team %s",
			c.Args[0].Value, team.Name)
		return fmt.Sprintf("Failed to add member %s to team %s",
			c.Args[0].Value, team.Name), noParams
	}
	return toMention(member.SlackID) +
		" was added to `" + team.Name + "` team :tada:", noParams
}
