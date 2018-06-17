package core

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewSetTechLeadCmd returns an add tech lead command that toggles an existing
// user's tech lead status (this action can only be performed by admins)
func NewSetTechLeadCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "set-techlead",
		HelpText: "Toggle an existing user's tech lead status (admins only)",
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

// setTechLead toggles an existing user's tech lead status
func (core *Plugin) setTechLead(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}
	username := c.Options["user"].Value
	member := model.Member{
		SlackID: cmd.ParseMention(username),
		IsAdmin: true,
	}
	if err := core.Bot.DAL.SetMemberIsTechLead(&member); err != nil {
		log.WithError(err).Error("Failed to update %s's tech lead status", username)
		return "Failed to update admin status", noParams
	}
	return cmd.ToMention(member.SlackID) + "'s tech lead status has been updated :tada:", noParams
}
