package bot

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
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:      "username",
				HelpText:  "the Slack handle of the user to remove admin rights from",
				Format:    anyRegex,
				MultiWord: false,
			},
		},
		HandleFunc: ch,
	}
}

// removeAdmin removes admin priveledges from an existing user.
func (b *Bot) removeAdmin(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}
	user := model.Member{
		SlackID: parseMention(c.Args[0].Value),
		IsAdmin: false,
	}
	if err := b.dal.SetMemberIsAdmin(&user); err != nil {
		return "Failed to remove user's admin priveleges", noParams
	}
	return toMention(user.SlackID) + " has been removed as admin :tada:", noParams
}
