package cmd

import (
	"github.com/nlopes/slack"
	"github.com/ubclaunchpad/rocket/model"
)

// Context stores a Slack message and the user who sent it.
type Context struct {
	Message *slack.Msg
	User    model.Member
	Options map[string]Option
}

// CommandHandler is the interface all handlers of Rocket commands must implement.
type CommandHandler func(Context) (string, slack.PostMessageParameters)
