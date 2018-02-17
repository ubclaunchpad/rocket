package bot

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// Send a help message
func (b *Bot) help(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	res := ""
	opt := c.Options["command"].Value
	if opt == "" {
		// General help
		cmds := ""
		for _, cmd := range Commands {
			cmds += fmt.Sprintf("\t%s\t\t%s\n", cmd.Name, cmd.HelpText)
		}
		res = fmt.Sprintf("Usage: @rocket COMMAND\n\nGet help using a specific "+
			"command with \"@rocket help --command=`COMMAND`\"\n\nCommands:\n%s",
			cmds)
		return res, noParams
	}
	// Command-specific help
	for _, cmd := range Commands {
		if opt == cmd.Name {
			return cmd.Help(), noParams
		}
	}
	res = fmt.Sprintf("\"%s\" is not a Rocket command.\n"+
		"See \"@rocket help\"", opt)
	return res, noParams
}

// Generic command for setting some information about the sender's profile.
func (b *Bot) set(c cmd.Context) (string, slack.PostMessageParameters) {
	params := slack.PostMessageParameters{}

	if c.Options["name"].Value != "" {
		c.User.Name = c.Options["name"].Value
		if err := b.dal.SetMemberName(&c.User); err != nil {
			log.WithError(err).Errorf("Failed to set name: %s", c.User.Name)
			return "Failed to set name " + c.User.Name, params
		}
	}

	if c.Options["email"].Value != "" {
		c.User.Email = c.Options["email"].Value
		if err := b.dal.SetMemberEmail(&c.User); err != nil {
			log.WithError(err).Errorf("Failed to set email: %s", c.User.Email)
			return "Failed to set email " + c.User.Email, params
		}
	}

	if c.Options["github"].Value != "" {
		c.User.GithubUsername = c.Options["github"].Value
		// Check that the user exists
		exists, err := b.gh.UserExists(c.User.GithubUsername)
		if err != nil {
			log.WithError(err).Errorf("Error checking whether user %s exists", c.User.GithubUsername)
			return "Error checking whether user exists", params
		} else if !exists {
			return fmt.Sprintf("Github user %s does not exist", c.User.GithubUsername), params
		}

		// Add the user to our GitHub org by adding to `all` team
		if err := b.gh.AddUserToTeam(c.User.GithubUsername, githubAllTeamID); err != nil {
			log.WithError(err).Errorf("Failed to add %s to Launch Pad Github organization",
				c.User.GithubUsername)
			return "Failed to add you to Launch Pad's GitHub organization", params
		}

		// Finally, set their username in the DB
		if err := b.dal.SetMemberGitHubUsername(&c.User); err != nil {
			log.WithError(err).Errorf("Failed to set GitHub username")
			return "Failed to set GitHub username", params
		}
	}

	if c.Options["major"].Value != "" {
		c.User.Major = c.Options["major"].Value
		if err := b.dal.SetMemberMajor(&c.User); err != nil {
			log.WithError(err).Error("Failed to set major")
			return "Failed to set major", params
		}
	}

	if c.Options["position"].Value != "" {
		c.User.Position = c.Options["position"].Value
		if err := b.dal.SetMemberPosition(&c.User); err != nil {
			log.WithError(err).Error("Failed to set position")
			return "Failed to set position", params
		}
	}

	params.Attachments = c.User.SlackAttachments()
	return "Your position has been updated :simple_smile:", params
}

// addUser adds an existing user to a team
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
	if err := b.gh.AddUserToTeam(member.GithubUsername, team.GithubTeamID); err != nil {
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

// addTeam creates a new Launch Pad team
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

// addAdmin makes an existing user and admin
func (b *Bot) addAdmin(c cmd.Context) (string, slack.PostMessageParameters) {
	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}

	noParams := slack.PostMessageParameters{}
	user := model.Member{
		SlackID: parseMention(c.Args[0].Value),
		IsAdmin: true,
	}
	if err := b.dal.SetMemberIsAdmin(&user); err != nil {
		log.WithError(err).Error("Failed to make user " + c.Args[0].Value + " admin")
		return "Failed to make user admin", noParams
	}
	return toMention(user.SlackID) + " has been made an admin :tada:", noParams
}

// removeTeam removes a Launch Pad team
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

// removeAdmin removes admin priveledges from an existing user
func (b *Bot) removeAdmin(c cmd.Context) (string, slack.PostMessageParameters) {
	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}

	noParams := slack.PostMessageParameters{}
	user := model.Member{
		SlackID: parseMention(c.Args[0].Value),
		IsAdmin: false,
	}
	if err := b.dal.SetMemberIsAdmin(&user); err != nil {
		return "Failed to remove user's admin priveleges", noParams
	}
	return toMention(user.SlackID) + " has been removed as admin :tada:", noParams
}

// removeUser removes a user from the Launch Pad organization
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
		log.WithError(err).Errorf("Failed to add member %s to GitHub team %s",
			c.Args[0].Value, team.Name)
		return "Failed to add member to GitHub team", noParams
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
		" was removed from `" + team.Name + "` team :tada:", noParams
}

// viewUser displays a user's information
func (b *Bot) viewUser(c cmd.Context) (string, slack.PostMessageParameters) {
	params := slack.PostMessageParameters{}
	user := model.Member{
		SlackID: parseMention(c.Args[0].Value),
	}
	if err := b.dal.GetMemberBySlackID(&user); err != nil {
		log.WithError(err).Error("Failed to get member " + c.Args[0].Value)
		return "Failed to get member " + c.Args[0].Value, params
	}
	params.Attachments = user.SlackAttachments()
	return c.Args[0].Value + "'s profile", params
}

// viewTeam displays a teams's information
func (b *Bot) viewTeam(c cmd.Context) (string, slack.PostMessageParameters) {
	params := slack.PostMessageParameters{}
	team := model.Team{
		Name: c.Args[0].Value,
	}
	if err := b.dal.GetTeamByName(&team); err != nil {
		log.WithError(err).Error("Failed to get team " + team.Name)
		return "Failed to get team " + team.Name, params
	}
	params.Attachments = team.SlackAttachments()
	return "Team " + c.Args[0].Value, params
}

// cmd.Command for debugging strange behaviour without restarting the whole app.
// It refreshes the user cache and creates any users that don't already exist.
func (b *Bot) refresh(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}

	// Pull in all users from Slack
	b.PopulateUsers()

	var member model.Member
	for _, user := range b.users {
		member = model.Member{
			SlackID:  user.ID,
			ImageURL: user.Profile.Image192,
		}

		if err := b.dal.CreateMember(&member); err != nil {
			log.WithError(err).Error("Error creating member with Slack ID " + member.SlackID)
			return "Error creating member with Slack ID " + member.SlackID, noParams
		}

		// Set Slack image URL
		if err := b.dal.SetMemberImageURL(&member); err != nil {
			b.log.WithError(err).Error("Error setting image for Slack ID " + member.SlackID)
			return "Error setting image for Slack ID %s" + member.SlackID, noParams
		}
	}
	return "I feel so refreshed! :tropical_drink:", noParams
}
