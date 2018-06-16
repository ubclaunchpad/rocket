package core

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewAddTechLeadCmd returns an add tech leadcommand that makes an existing user
// a tech lead (this action can only be performed by admins)
func NewAddTechLeadCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "add-techlead",
		HelpText: "Make an existing user a tech lead (admins only)",
		Options: map[string]*cmd.Option{
			"user": &cmd.Option{
				Key:      "user",
				HelpText: "the Slack handle of the user to make a tech lead",
				Format:   cmd.AnyRegex,
				Required: true,
			},
		},
		HandleFunc: ch,
	}
}

// addTechLead makes an existing user and turns them into a tech lead
func (core *CorePlugin) addTechLead(c cmd.Context) (string, slack.PostMessageParameters) {
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
		log.WithError(err).Error("Failed to make user " + username + " tech lead")
		return "Failed to make user tech lead", noParams
	}
	return cmd.ToMention(member.SlackID) + " has been made a tech lead :tada:", noParams
}
