package bot

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewAddAdminCmd returns an add admin command that makes an existing user an
// admin (this action can only be performed by admins)
func NewAddAdminCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "add-admin",
		HelpText: "Make an existing user an admin (admins only)",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:      "username",
				HelpText:  "the Slack handle of the user to make an admin",
				Format:    anyRegex,
				MultiWord: false,
			},
		},
		HandleFunc: ch,
	}
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
