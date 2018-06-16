package core

import (
	"github.com/nlopes/slack"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewRemoveAdminCmd returns a remove admin command that makes an existing
// user an admin (this action can only be
// performed by admins)
func NewRemoveAdminCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "remove-admin",
		HelpText: "Remove admin rights from a user (admins only)",
		Options: map[string]*cmd.Option{
			"user": &cmd.Option{
				Key:      "user",
				HelpText: "the Slack handle of the user to remove admin rights from",
				Format:   cmd.AnyRegex,
				Required: true,
			},
		},
		HandleFunc: ch,
	}
}

// removeAdmin removes admin priveledges from an existing user.
func (core *Plugin) removeAdmin(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}
	user := model.Member{
		SlackID: cmd.ParseMention(c.Options["user"].Value),
		IsAdmin: false,
	}
	if err := core.Bot.DAL.SetMemberIsAdmin(&user); err != nil {
		return "Failed to remove user's admin priveleges", noParams
	}
	return cmd.ToMention(user.SlackID) + " has been removed as admin :tada:", noParams
}
