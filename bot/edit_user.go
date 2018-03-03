package bot

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
			"name": &cmd.Option{
				Key:      "name",
				HelpText: "user's full name",
				Format:   nameRegex,
			},
			"email": &cmd.Option{
				Key:      "email",
				HelpText: "user's email address",
				Format:   emailRegex,
			},
			"position": &cmd.Option{
				Key:      "position",
				HelpText: "user's creative Launch Pad title",
				Format:   anyRegex,
			},
			"github": &cmd.Option{
				Key:      "github",
				HelpText: "user's Github username",
				Format:   anyRegex,
			},
			"major": &cmd.Option{
				Key:      "major",
				HelpText: "user's major at UBC",
				Format:   anyRegex,
			},
		},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:      "member",
				HelpText:  "the Slack handle of the user to edit",
				Format:    anyRegex,
				MultiWord: false,
			},
		},
		HandleFunc: ch,
	}
}

// Generic command for setting some information about the sender's profile.
func (b *Bot) editUser(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}

	c.User = model.Member{
		SlackID: parseMention(c.Args[0].Value),
	}
	if err := b.dal.GetMemberBySlackID(&c.User); err != nil {
		return "Failed to find member " + c.Args[0].Value, noParams
	}
	_, params := b.set(c)
	return c.Args[0].Value + "'s information has been updated", params
}
