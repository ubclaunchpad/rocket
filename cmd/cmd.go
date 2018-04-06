package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	"github.com/ubclaunchpad/rocket/model"
)

// Command represents a command that Rocket will recognise and respond to.
type Command struct {
	Name       string
	HelpText   string
	Options    map[string]*Option
	HandleFunc CommandHandler
}

// Execute executes the given command and returns an error if
// if the command is invalid.
func (c *Command) Execute(ctx Context) (string, slack.PostMessageParameters, error) {
	// Parse and validate command
	if err := c.parse(ctx.Message.Text); err != nil {
		return "", slack.PostMessageParameters{}, err
	}
	// Copy options to context for use by command handler
	ctx.Options = map[string]Option{}
	for key, opt := range c.Options {
		ctx.Options[key] = *opt
		// Clear option value now that it's been copied
		opt.Value = ""
	}
	// Pass context to command handler
	res, params := c.HandleFunc(ctx)
	return res, params, nil
}

// Help returns full help text for the given command
func (c *Command) Help() (string, slack.PostMessageParameters) {
	usage := "Usage: @rocket " + c.Name
	opts := ""
	attachments := []slack.Attachment{}
	if len(c.Options) > 0 {
		usage += " OPTIONS"
		opts = ""
		for _, o := range c.Options {
			if o.Required {
				opts += fmt.Sprintf("%s (required): %s\n", o.Key, o.HelpText)
			} else {
				opts += fmt.Sprintf("%s: %s\n", o.Key, o.HelpText)
			}
		}
		attachments = append(attachments, slack.Attachment{
			Title: "Options",
			Text:  opts,
			Color: "#e5e7ea",
		})
	}
	params := slack.PostMessageParameters{Attachments: attachments}
	return fmt.Sprintf("%s\n\n%s", usage, c.HelpText), params
}

// parse checks whether the given command meets the requirements of this
// Command and returns nil if it does, and the validation error otherwise.
// The command format should be "@rocket COMMAND OPTIONS ARGUMENTS"
func (c *Command) parse(cmd string) error {
	// Check that we received the correct command
	tokens := strings.Fields(cmd)
	if len(tokens) < 2 {
		return fmt.Errorf("Received empty command")
	} else if tokens[1] != c.Name {
		return fmt.Errorf("Invalid command \"%s\"", tokens[1])
	}
	// Check options and store their values
	optionsRegex := regexp.MustCompile("[a-zA-Z-]+=`[^`]+`")
	opts := optionsRegex.FindAllString(strings.Join(tokens[2:], " "), -1)
	return c.parseOptions(opts)
}

// parseOptions checks that the value corresponding to each option matches
// that option's required format, then stores that value. Returns an error
// if an option is malformatted, or a required option is missing.
// opts should be a slice of strings of the format "key=value".
func (c *Command) parseOptions(opts []string) error {
	for _, token := range opts {
		// Token has format my-key=`my value`. Extract option key and value
		parts := strings.SplitN(token, "=", 2)
		key := parts[0][:len(parts[0])]
		value := strings.TrimRight(strings.TrimLeft(parts[1], "`"), "`")

		// Check that it is a valid option
		option := c.Options[key]
		if option == nil {
			return fmt.Errorf("Unrecognized option %s", key)
		}

		// Check that the option fits it's specified format
		if err := option.validate(value); err != nil {
			return err
		}

		// Check that this option has not already been set
		if option.Value != "" {
			return fmt.Errorf("Duplicate option \"%s\"", key)
		}
		option.Value = value
	}

	// Check that we aren't missing any required options
	for _, option := range c.Options {
		if option.Required && option.Value == "" {
			return fmt.Errorf("Missing value for required option \"%s\"", option.Key)
		}
	}
	return nil
}

// Option represents a parameter that can be passed as part of a
// Rocket command
type Option struct {
	Key      string
	HelpText string
	Format   *regexp.Regexp
	Required bool
	Value    string
}

// validate returns nil if the given value meets the format requirements for
// this option, returns the validation error otherwise.
func (o *Option) validate(value string) error {
	// Check that the value meets the required format
	if !o.Format.MatchString(value) {
		return fmt.Errorf("Invalid format for option \"%s\". "+
			"Format must match regular expression %s.", o.Key, o.Format.String())
	}
	return nil
}

// Context stores a Slack message and the user who sent it.
type Context struct {
	Message *slack.Msg
	User    model.Member
	Options map[string]Option
}

// CommandHandler is the interface all handlers of Rocket commands must implement.
type CommandHandler func(Context) (string, slack.PostMessageParameters)
