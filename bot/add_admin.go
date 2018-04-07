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
		Options: map[string]*cmd.Option{
			"user": &cmd.Option{
				Key:      "user",
				HelpText: "the Slack handle of the user to make an admin",
				Format:   anyRegex,
				Required: true,
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
	username := c.Options["user"].Value
	member := model.Member{
		SlackID: parseMention(username),
		IsAdmin: true,
	}
	if err := b.dal.SetMemberIsAdmin(&member); err != nil {
		log.WithError(err).Error("Failed to make user " + username + " admin")
		return "Failed to make user admin", noParams
	}
	return toMention(member.SlackID) + " has been made an admin :tada:", noParams
}
