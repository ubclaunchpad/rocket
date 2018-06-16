package core

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
				Format:   cmd.AnyRegex,
				Required: true,
			},
		},
		HandleFunc: ch,
	}
}

// addAdmin makes an existing user and admin
func (core *RocketPlugin) addAdmin(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}
	username := c.Options["user"].Value
	member := model.Member{
		SlackID: cmd.ParseMention(username),
		IsAdmin: true,
	}
	if err := core.Bot.DAL.SetMemberIsAdmin(&member); err != nil {
		log.WithError(err).Error("Failed to make user " + username + " admin")
		return "Failed to make user admin", noParams
	}
	return cmd.ToMention(member.SlackID) + " has been made an admin :tada:", noParams
}
