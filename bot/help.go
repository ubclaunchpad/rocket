package bot

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/ubclaunchpad/rocket/cmd"
)

// NewHelpCmd retuns a help command that presents helpful information about
// Rocket commands.
func NewHelpCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "help",
		HelpText: "Get help using Rocket commands",
		Options: map[string]*cmd.Option{
			"command": &cmd.Option{
				Key:      "command",
				HelpText: "get help using a particular Rocket command",
				Format:   alphaRegex,
				Required: false,
			},
		},
		HandleFunc: ch,
	}
}

// Send a help message
func (b *Bot) help(c cmd.Context) (string, slack.PostMessageParameters) {
	params := slack.PostMessageParameters{}
	res := ""
	opt := c.Options["command"].Value
	if opt == "" {
		// General help
		res = "Usage: @rocket COMMAND\n\nGet help using a specific " +
			"command with \"@rocket help command={COMMAND}\"\n" +
			"Example: @rocket set name={A Guy} github={arealguy}"
		cmds := ""
		for _, cmd := range b.commands {
			cmds += fmt.Sprintf("%s\t\t%s\n", cmd.Name, cmd.HelpText)
		}
		commands := slack.Attachment{
			Title: "Commands",
			Text:  cmds,
			Color: "#e5e7ea",
		}
		params.Attachments = []slack.Attachment{commands}
		return res, params
	}
	// Command-specific help
	for _, cmd := range b.commands {
		if opt == cmd.Name {
			return cmd.Help()
		}
	}
	res = fmt.Sprintf("\"%s\" is not a Rocket command.\n"+
		"See \"@rocket help\"", opt)
	return res, noParams
}
