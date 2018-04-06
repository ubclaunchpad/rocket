package bot

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewViewUserCmd returns a view user command that displays information about a user
func NewViewUserCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "view-user",
		HelpText: "View information about a user",
		Options: map[string]*cmd.Option{
			"user": &cmd.Option{
				Key:      "user",
				HelpText: "the slack handle of the user to view",
				Format:   anyRegex,
				Required: true,
			},
		},
		HandleFunc: ch,
	}
}

// viewUser displays a user's information.
func (b *Bot) viewUser(c cmd.Context) (string, slack.PostMessageParameters) {
	params := slack.PostMessageParameters{}
	username := c.Options["user"].Value
	user := model.Member{
		SlackID: parseMention(username),
	}
	if err := b.dal.GetMemberBySlackID(&user); err != nil {
		log.WithError(err).Error("Failed to get member " + username)
		return "Failed to get member " + username, params
	}
	params.Attachments = user.SlackAttachments()
	return username + "'s profile", params
}
