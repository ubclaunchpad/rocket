package bot

import (
	"strings"

	"github.com/nlopes/slack"
)

func (b *Bot) help(c *CommandContext) {
	b.api.PostMessage(c.msg.Channel, helpMessage, noParams)
}

func (b *Bot) me(c *CommandContext) {
	params := slack.PostMessageParameters{}
	params.Attachments = c.user.SlackAttachments()
	b.api.PostMessage(c.msg.Channel, "Your Launch Pad profile :rocket:", params)
}

func (b *Bot) set(c *CommandContext) {
	params := slack.PostMessageParameters{}
	switch c.args[2] {
	case "name":
		c.user.Name = strings.Join(c.args[3:], " ")
		if err := b.dal.SetMemberName(&c.user); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to set name")
		}
		params.Attachments = c.user.SlackAttachments()
		b.api.PostMessage(c.msg.Channel, "Your name has been updated! :simple_smile:", params)
	case "email":
		c.user.Email = parseEmail(c.args[3])
		if err := b.dal.SetMemberEmail(&c.user); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to set email")
		}
		params.Attachments = c.user.SlackAttachments()
		b.api.PostMessage(c.msg.Channel, "You email has been updated :simple_smile:", params)
	case "github":
		c.user.GithubUsername = c.args[3]
		if err := b.dal.SetMemberGitHubUsername(&c.user); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to set github username")
		}
		params.Attachments = c.user.SlackAttachments()
		b.api.PostMessage(c.msg.Channel, "Your GitHub username has been updated :simple_smile:", params)
	case "major":
		c.user.Major = c.args[3]
		if err := b.dal.SetMemberMajor(&c.user); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to set major")
		}
		params.Attachments = c.user.SlackAttachments()
		b.api.PostMessage(c.msg.Channel, "Your major has bee updated :simple_smile:", params)
	case "position":
		c.user.Position = strings.Join(c.args[3:], " ")
		if err := b.dal.SetMemberPosition(&c.user); err != nil {
			b.SendErrorMessage(c.msg.Channel, err, "Failed to set ")
		}
		params.Attachments = c.user.SlackAttachments()
		b.api.PostMessage(c.msg.Channel, "Your position has been updated :simple_smile:", params)
	}
}
