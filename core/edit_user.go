package core

import (
	"github.com/nlopes/slack"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewEditUserCmd returns an edit user command that allows admins to edit other
// users' info
func NewEditUserCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "edit",
		HelpText: "Set properties on another user's Launch Pad profile (admins only)",
		Options: map[string]*cmd.Option{
			"member": &cmd.Option{
				Key:      "member",
				HelpText: "the Slack handle of the user to edit",
				Format:   cmd.AnyRegex,
				Required: false,
			},
			"name": &cmd.Option{
				Key:      "name",
				HelpText: "user's full name",
				Format:   cmd.NameRegex,
				Required: false,
			},
			"email": &cmd.Option{
				Key:      "email",
				HelpText: "user's email address",
				Format:   cmd.EmailRegex,
				Required: false,
			},
			"position": &cmd.Option{
				Key:      "position",
				HelpText: "user's creative Launch Pad title",
				Format:   cmd.AnyRegex,
				Required: false,
			},
			"github": &cmd.Option{
				Key:      "github",
				HelpText: "user's Github username",
				Format:   cmd.AnyRegex,
				Required: false,
			},
			"major": &cmd.Option{
				Key:      "major",
				HelpText: "user's major at UBC",
				Format:   cmd.AnyRegex,
				Required: false,
			},
		},
		HandleFunc: ch,
	}
}

// Generic command for setting some information about the sender's profile.
func (core *RocketPlugin) editUser(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}

	memberName := c.Options["member"].Value
	c.User = model.Member{
		SlackID: cmd.ParseMention(memberName),
	}
	if err := core.Bot.DAL.GetMemberBySlackID(&c.User); err != nil {
		return "Failed to find member " + memberName, noParams
	}
	_, params := core.set(c)
	return memberName + "'s information has been updated", params
}
