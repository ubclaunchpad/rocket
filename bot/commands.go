package bot

import (
	"context"
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
		b.api.PostMessage(c.msg.Channel, "You email has been updated :simple_smile:", params)
	case "github":
		c.user.GithubUsername = c.args[2]
		// Check that the user exists
		user, _, err := b.gh.Users.Get(context.Background(), c.args[2])
		if err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "GitHub user does not exist")
			return
		}
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
		team := model.Team{
			Name: strings.Join(c.args[2:], " "),
		}
		if err := b.dal.CreateTeam(&team); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to create team")
			return
		}
		b.api.PostMessage(c.msg.Channel, "`"+team.Name+"` team has been created :tada:", noParams)
	case "admin":
		user := model.Member{
			SlackID: parseMention(c.args[3]),
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

		member := model.TeamMember{
			MemberSlackID: parseMention(c.args[1]),
			TeamName:      c.args[3],
		}
		if err := b.dal.CreateTeamMember(&member); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to add member to team")
			return
		}
		b.api.PostMessage(c.msg.Channel, toMention(member.MemberSlackID)+" was added to `"+member.TeamName+"` team :tada:", noParams)
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
		if err := b.dal.DeleteTeam(&team); err != nil {
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

		member := model.TeamMember{
			MemberSlackID: parseMention(c.args[1]),
			TeamName:      c.args[3],
		}
		if err := b.dal.DeleteTeamMember(&member); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to remove member from team")
			return
		}
		b.api.PostMessage(c.msg.Channel, toMention(member.MemberSlackID)+" was removed from `"+member.TeamName+"` team :tada:", noParams)
	}
}

func (b *Bot) refresh(c *CommandContext) {
	b.PopulateUsers()
	b.api.PostMessage(c.msg.Channel, "I feel so refreshed! :tropical_drink:", noParams)
}
