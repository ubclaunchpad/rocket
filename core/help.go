package core

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
				Format:   cmd.AlphaRegex,
				Required: false,
			},
		},
		HandleFunc: ch,
	}
}

// Send a help message
func (core *Plugin) help(c cmd.Context) (string, slack.PostMessageParameters) {
	params := slack.PostMessageParameters{}
	res := ""
	opt := c.Options["command"].Value
	if opt == "" {
		// General help
		res = "Usage: `@rocket COMMAND`\n\nGet help using a specific " +
			"command with `@rocket help command={COMMAND}`\n" +
			"Example: `@rocket set name={A Guy} github={arealguy}`"

		// Get length of longest command to generate nicely spaced spaces
		// between command name and the corresponding description
		longest := 0
		for _, cmd := range core.Bot.Commands {
			if len(cmd.Name) > longest {
				longest = len(cmd.Name)
			}
		}
		maxDividerSpace := longest + 1

		// Format help text
		cmds := "```\n"
		for _, cmd := range core.Bot.Commands {
			dividerSpace := ""
			for i := 0; i < maxDividerSpace-len(cmd.Name); i++ {
				dividerSpace += " "
			}
			cmds += fmt.Sprintf("%s%s%s\n", cmd.Name, dividerSpace, cmd.HelpText)
		}
		cmds += "\n```"
		commands := slack.Attachment{
			Title: "Commands",
			Text:  cmds,
			Color: "#e5e7ea",
		}
		params.Attachments = []slack.Attachment{commands}
		return res, params
	}
	// Command-specific help
	for _, cmd := range core.Bot.Commands {
		if opt == cmd.Name {
			return cmd.Help()
		}
	}
	res = fmt.Sprintf("`%s` is not a Rocket command.\n"+
		"See `@rocket help`", opt)
	return res, slack.PostMessageParameters{}
}
