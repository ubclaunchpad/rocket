package core

import (
	"fmt"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewToggleAdminCmd returns an add admin command that makes an existing user an
// admin (this action can only be performed by admins)
func NewToggleAdminCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "toggle-admin",
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

// toggleAdmin toggles an existing user's admin status
func (core *Plugin) toggleAdmin(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}
	username := c.Options["user"].Value
	member := &model.Member{SlackID: cmd.ParseMention(username)}
	err := core.Bot.DAL.GetMemberBySlackID(member)
	if err != nil {
		log.WithError(err).Error("Failed to get %s", username)
		return "Failed to find user", noParams
	}

	// Update member admin status
	member.IsAdmin = !member.IsAdmin
	if err := core.Bot.DAL.SetMemberIsAdmin(member); err != nil {
		log.WithError(err).Error("Failed to update %s's admin status", username)
		return "Failed to update admin status", noParams
	}
	return fmt.Sprintf(
		"Set %s's admin status has been set to %t :tada:",
		cmd.ToMention(member.SlackID), member.IsAdmin,
	), noParams
}
