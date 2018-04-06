package bot

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewRefreshCmd returns a refresh command that refreshes the user cache and creates
// any users that don't already exist
func NewRefreshCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "refresh",
		HelpText: "for debugging Rocket (admins only)",
		Options:  map[string]*cmd.Option{},
	}
}

// refresh is a command for debugging strange behaviour without restarting the
// whole app. It refreshes the user cache and creates any users that don't
// already exist.
func (b *Bot) refresh(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	if !c.User.IsAdmin {
		return "You must be an admin to use this command", noParams
	}

	// Pull in all users from Slack
	b.PopulateUsers()

	var member model.Member
	for _, user := range b.users {
		member = model.Member{
			SlackID:  user.ID,
			ImageURL: user.Profile.Image192,
		}

		if err := b.dal.CreateMember(&member); err != nil {
			log.WithError(err).Error("Error creating member with Slack ID " + member.SlackID)
			return "Error creating member with Slack ID " + member.SlackID, noParams
		}

		// Set Slack image URL
		if err := b.dal.SetMemberImageURL(&member); err != nil {
			b.log.WithError(err).Error("Error setting image for Slack ID " + member.SlackID)
			return "Error setting image for Slack ID %s" + member.SlackID, noParams
		}
	}
	return "I feel so refreshed! :tropical_drink:", noParams
}
