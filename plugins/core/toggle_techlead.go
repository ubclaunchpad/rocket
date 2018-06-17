package core

import (
	"fmt"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewToggleTechLeadCmd returns an add tech lead command that toggles an existing
// user's tech lead status (this action can only be performed by admins)
func NewToggleTechLeadCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "toggle-techlead",
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

// toggleTechLead toggles an existing user's tech lead status
func (core *Plugin) toggleTechLead(c cmd.Context) (string, slack.PostMessageParameters) {
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

	// Update tech lead status
	member.IsTechLead = !member.IsTechLead
	if err := core.Bot.DAL.SetMemberIsTechLead(member); err != nil {
		log.WithError(err).Error("Failed to update %s's tech lead status", username)
		return "Failed to update tech lead status", noParams
	}
	return fmt.Sprintf(
		"Set %s's tech lead status has been set to %t :tada:",
		cmd.ToMention(member.SlackID), member.IsTechLead,
	), noParams
}
