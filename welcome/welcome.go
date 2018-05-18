package welcome

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/ubclaunchpad/rocket/bot"
	"github.com/ubclaunchpad/rocket/cmd"
)

// WelcomePlugin stores the bot that is used to access the Slack API.
type WelcomePlugin struct {
	Bot *bot.Bot
}

// New reutrns a new instance of the WelcomePlugin
func New(b *bot.Bot) *WelcomePlugin {
	return &WelcomePlugin{
		Bot: b,
	}
}

// Start starts the welcome plugin.
func (wp *WelcomePlugin) Start() error {
	wp.Bot.Log.Info("Running WelcomePlugin")
	return nil
}

// Commands returns an empty list of commands, because this plugin has no
// commands.
func (wp *WelcomePlugin) Commands() []*cmd.Command {
	return []*cmd.Command{}
}

// EventHandlers returns a map from event type to event handler.
func (wp *WelcomePlugin) EventHandlers() map[string]bot.EventHandler {
	return map[string]bot.EventHandler{
		"team_join": wp.handleTeamJoin,
	}
}

// handleTeamJoin welcomes a user to our Slack when they join be messaging
// them in the general channel.
func (wp *WelcomePlugin) handleTeamJoin(evt slack.RTMEvent) {
	user := evt.Data.(*slack.TeamJoinEvent).User
	msg := fmt.Sprintf("Welcome to the team, %s! :rocket:", cmd.ToMention(user.ID))
	noParams := slack.PostMessageParameters{}
	wp.Bot.API.PostMessage("general", msg, noParams)
}
