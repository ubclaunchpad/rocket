package core

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewSetAdminCmd returns an add admin command that makes an existing user an
// admin (this action can only be performed by admins)
func NewSetAdminCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "set-admin",
		HelpText: "Toggle an existing user's admin status (admins only)",
		Options: map[string]*cmd.Option{
			"user": &cmd.Option{
				Key:      "user",
				HelpText: "the Slack handle of the user to update",
				Format:   cmd.AnyRegex,
				Required: true,
			},
		},
		HandleFunc: ch,
	}
}

// setAdmin toggles an existing user's admin status
func (core *Plugin) setAdmin(c cmd.Context) (string, slack.PostMessageParameters) {
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
		log.WithError(err).Error("Failed to update %s's admin status", username)
		return "Failed to update admin status", noParams
	}
	return cmd.ToMention(member.SlackID) + "'s admin status has been updated :tada:", noParams
}
