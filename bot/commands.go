package bot

import (
	"strconv"
	"strings"

	"github.com/nlopes/slack"
	"github.com/ubclaunchpad/rocket/model"
)

// Command handlers accept a string slice of the form
// <command name> <arg1> <arg2> ... <argN>

func (b *Bot) help(c *CommandContext) {
	b.api.PostMessage(c.msg.Channel, helpMessage, noParams)
}

func (b *Bot) me(c *CommandContext) {
	params := slack.PostMessageParameters{}
	params.Attachments = c.user.SlackAttachments()
	b.api.PostMessage(c.msg.Channel, "Your Launch Pad profile :rocket:", params)
}

func (b *Bot) set(c *CommandContext) {
	if len(c.args) < 3 {
		b.SendErrorMessage(c.msg.Channel, nil, "Not enough arguments")
		return
	}
	params := slack.PostMessageParameters{}
	switch c.args[1] {
	case "name":
		c.user.Name = strings.Join(c.args[2:], " ")
		if err := b.dal.SetMemberName(&c.user); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to set name")
			return
		}
		params.Attachments = c.user.SlackAttachments()
		b.api.PostMessage(c.msg.Channel, "Your name has been updated! :simple_smile:", params)
	case "email":
		c.user.Email = parseEmail(c.args[2])
		if err := b.dal.SetMemberEmail(&c.user); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to set email")
			return
		}
		params.Attachments = c.user.SlackAttachments()
		b.api.PostMessage(c.msg.Channel, "Your email has been updated :simple_smile:", params)
	case "github":
		c.user.GithubUsername = c.args[2]
		// Check that the user exists
		exists, err := b.gh.UserExists(c.user.GithubUsername)
		if err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Error checking whether user exists")
			return
		}
		if !exists {
			b.SendErrorMessage(c.msg.Channel, nil, "No GitHub user with that name exists")
			return
		}

		// Add the user to our GitHub org by adding to `all` team
		if err := b.gh.AddUserToTeam(c.user.GithubUsername, githubAllTeamID); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to add you to Launch Pad's GitHub organization")
			return
		}

		// Finally, set their username in the DB
		if err := b.dal.SetMemberGitHubUsername(&c.user); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to set github username")
			return
		}

		params.Attachments = c.user.SlackAttachments()
		b.api.PostMessage(c.msg.Channel, "Your GitHub username has been updated :simple_smile:", params)
	case "major":
		c.user.Major = c.args[2]
		if err := b.dal.SetMemberMajor(&c.user); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to set major")
			return
		}
		params.Attachments = c.user.SlackAttachments()
		b.api.PostMessage(c.msg.Channel, "Your major has been updated :simple_smile:", params)
	case "position":
		c.user.Position = strings.Join(c.args[2:], " ")
		if err := b.dal.SetMemberPosition(&c.user); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to set position")
			return
		}
		params.Attachments = c.user.SlackAttachments()
		b.api.PostMessage(c.msg.Channel, "Your position has been updated :simple_smile:", params)
	}
}

func (b *Bot) add(c *CommandContext) {
	if len(c.args) < 3 {
		b.SendErrorMessage(c.msg.Channel, nil, "Not enough arguments")
		return
	}

	if !c.user.IsAdmin {
		b.SendErrorMessage(c.msg.Channel, nil, "You must be an admin to use this command")
		return
	}

	switch c.args[1] {
	case "team":
		teamName := strings.Join(c.args[2:], " ")
		// teamName = "Great Team", ghTeamName = "great-team"
		ghTeamName := strings.ToLower(strings.Join(c.args[2:], "-"))

		// Create the team on GitHub
		ghTeam, err := b.gh.CreateTeam(ghTeamName)
		b.log.Info("create team, ", ghTeam, err)
		if err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to create team on GitHub")
			return
		}

		team := model.Team{
			Name:         teamName,
			GithubTeamID: *ghTeam.ID,
		}
		// Finally, add team to DB
		if err := b.dal.CreateTeam(&team); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to create team")
			return
		}
		b.api.PostMessage(c.msg.Channel, "`"+team.Name+"` has been added :tada:", noParams)
	case "admin":
		user := model.Member{
			SlackID: parseMention(c.args[2]),
			IsAdmin: true,
		}
		if err := b.dal.SetMemberIsAdmin(&user); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to make user admin")
			return
		}
		b.api.PostMessage(c.msg.Channel, toMention(user.SlackID)+" has been made an admin :tada:", noParams)
	default:
		if len(c.args) < 4 {
			b.SendErrorMessage(c.msg.Channel, nil, "Not enough arguments")
			return
		}

		team := model.Team{
			Name: strings.Join(c.args[3:], " "),
		}
		if err := b.dal.GetTeamByName(&team); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to find team")
			return
		}

		member := model.Member{
			SlackID: parseMention(c.args[1]),
		}
		if err := b.dal.GetMemberBySlackID(&member); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to find member")
			return
		}

		// Add user to corresponding GitHub team
		if err := b.gh.AddUserToTeam(member.GithubUsername, team.GithubTeamID); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to add user to GitHub team")
			return
		}

		teamMember := model.TeamMember{
			MemberSlackID: parseMention(c.args[1]),
			GithubTeamID:  team.GithubTeamID,
		}
		// Finally, add relation to DB
		if err := b.dal.CreateTeamMember(&teamMember); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to add member to team")
			return
		}
		b.api.PostMessage(c.msg.Channel, toMention(member.SlackID)+" was added to `"+team.Name+"` team :tada:", noParams)
	}
}

func (b *Bot) remove(c *CommandContext) {
	if len(c.args) < 3 {
		b.SendErrorMessage(c.msg.Channel, nil, "Not enough arguments")
		return
	}

	if !c.user.IsAdmin {
		b.SendErrorMessage(c.msg.Channel, nil, "You must be an admin to use this command")
		return
	}

	switch c.args[1] {
	case "team":
		team := model.Team{
			Name: strings.Join(c.args[2:], " "),
		}
		if err := b.dal.GetTeamByName(&team); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to find team")
			return
		}

		// Remove team from GitHub
		if err := b.gh.RemoveTeam(team.GithubTeamID); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to remove GitHub team")
			return
		}

		// Finally remove team from database
		if err := b.dal.DeleteTeamByName(&team); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to delete team")
			return
		}
		b.api.PostMessage(c.msg.Channel, "`"+team.Name+"` team has been deleted :tada:", noParams)
	case "admin":
		user := model.Member{
			SlackID: parseMention(c.args[2]),
			IsAdmin: false,
		}
		if err := b.dal.SetMemberIsAdmin(&user); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to remove user's admin priveleges")
			return
		}
		b.api.PostMessage(c.msg.Channel, toMention(user.SlackID)+" has been removed as admin :tada:", noParams)
	default:
		if len(c.args) < 4 {
			b.SendErrorMessage(c.msg.Channel, nil, "Not enough arguments")
			return
		}

		team := model.Team{
			Name: c.args[3],
		}
		if err := b.dal.GetTeamByName(&team); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to get team")
			return
		}

		member := model.Member{
			SlackID: parseMention(c.args[1]),
		}
		if err := b.dal.GetMemberBySlackID(&member); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to get member")
			return
		}

		// Remove user from GitHub team
		if err := b.gh.RemoveUserFromTeam(member.GithubUsername, team.GithubTeamID); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to add member to GitHub team")
			return
		}

		teamMember := model.TeamMember{
			MemberSlackID: parseMention(c.args[1]),
			GithubTeamID:  team.GithubTeamID,
		}
		// Remove user team relation from DB
		if err := b.dal.DeleteTeamMember(&teamMember); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to remove member from team")
			return
		}
		b.api.PostMessage(c.msg.Channel, toMention(member.SlackID)+" was removed from `"+team.Name+"` team :tada:", noParams)
	}
}

func (b *Bot) view(c *CommandContext) {
	if len(c.args) < 3 {
		b.SendErrorMessage(c.msg.Channel, nil, "Not enough arguments")
		return
	}

	switch c.args[1] {
	case "user":
		user := model.Member{
			SlackID: parseMention(c.args[2]),
		}
		if err := b.dal.GetMemberBySlackID(&user); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to get member")
			return
		}
		params := slack.PostMessageParameters{}
		params.Attachments = user.SlackAttachments()
		b.api.PostMessage(c.msg.Channel, c.args[2]+"'s profile", params)
	case "team":
		team := model.Team{
			Name: strings.Join(c.args[2:], " "),
		}
		if err := b.dal.GetTeamByName(&team); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to get team")
			return
		}
		params := slack.PostMessageParameters{}
		params.Attachments = team.SlackAttachments()
		b.api.PostMessage(c.msg.Channel, c.args[2]+" team", params)
	}
}

func (b *Bot) refresh(c *CommandContext) {
	// Pull in all users from Slack
	b.PopulateUsers()

	errCount := 0
	var member model.Member
	for _, user := range b.users {
		member = model.Member{
			SlackID:  user.ID,
			ImageURL: user.Profile.Image192,
		}

		if err := b.dal.CreateMember(&member); err != nil {
			b.log.WithError(err).Errorf("Error creating member with Slack ID %s", member.SlackID)
		}

		// Set Slack image URL
		if err := b.dal.SetMemberImageURL(&member); err != nil {
			b.log.WithError(err).Errorf("Error setting image for Slack ID %s", member.SlackID)
		}
	}

	if errCount > 0 {
		b.api.PostMessage(c.msg.Channel, strconv.Itoa(errCount)+" errors occurred while refreshing", noParams)
	}

	b.api.PostMessage(c.msg.Channel, "I feel so refreshed! :tropical_drink:", noParams)
}
