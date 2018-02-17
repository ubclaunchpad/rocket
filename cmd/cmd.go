package cmd

import (
	"errors"
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
	Args       []Argument
	HandleFunc CommandHandler
}

// Execute executes the given command and returns an error if
// if the command is invalid.
func (c *Command) Execute(ctx Context) (string, slack.PostMessageParameters, error) {
	// Parse and validate command
	if err := c.parse(ctx.Message.Text); err != nil {
		return "", slack.PostMessageParameters{}, err
	}
	// Copy options and args to context for use by command handler
	ctx.Options = map[string]Option{}
	for key, opt := range c.Options {
		ctx.Options[key] = *opt
		opt.Value = ""
	}
	ctx.Args = []Argument{}
	for _, arg := range c.Args {
		ctx.Args = append(ctx.Args, arg)
		arg.Value = ""
	}
	// Pass context to command handler
	res, params := c.HandleFunc(ctx)
	return res, params, nil
}

// Help returns full help text for the given command
func (c *Command) Help() string {
	usage := "Usage: @rocket " + c.Name
	opts := ""
	args := ""
	if len(c.Options) > 0 {
		usage += " OPTIONS"
		opts = "\nOptions:\n"
		for _, o := range c.Options {
			opts += fmt.Sprintf("\t--%s\t%s\n", o.Key, o.HelpText)
		}
	}
	if len(c.Args) > 0 {
		usage += " ARGUMENTS"
		args = "\nArguments:\n"
		for _, a := range c.Args {
			args += fmt.Sprintf("\t%s\t%s\n", a.Name, a.HelpText)
		}
	}
	return fmt.Sprintf("%s\n\n%s\n%s\n%s", usage, c.HelpText, args, opts)
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
		return fmt.Errorf("Invalid command \"%s\"", tokens[0])
	}
	if len(tokens) == 2 {
		// No options or args were given
		if len(c.Args) == 0 {
			return nil
		}
		return fmt.Errorf("Expected %d argument(s), but received 0", len(c.Args))
	}

	tokens = tokens[2:]
	argsAndOpts := strings.Join(tokens, " ")

	// Handle options
	optionsOnly := regexp.MustCompile("--[a-zA-Z0-9-]+=`[^`]+`")
	opts := optionsOnly.FindAllString(argsAndOpts, -1)
	if err := c.parseOptions(opts); err != nil {
		return err
	}

	// Remove options
	for _, o := range opts {
		argsAndOpts = strings.Replace(argsAndOpts, o, "", -1)
	}

	// Handle arguments
	args := strings.Fields(argsAndOpts)
	if err := c.parseArgs(args); err != nil {
		return err
	}
	return nil
}

func (c *Command) parseOptions(opts []string) error {
	for _, token := range opts {
		token = token[2:]

		// Extract option key and value
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
	return nil
}

func (c *Command) parseArgs(args []string) error {
	argIndex := 0
	for _, token := range args {
		if argIndex >= len(c.Args) {
			return errors.New("Too many arguments")
		}

		// Check that this argument fits it's specified format
		if err := c.Args[argIndex].validate(token); err != nil {
			return err
		}
		c.Args[argIndex].Value = token
		argIndex++
	}

	// Check that all args were provided
	if argIndex != len(c.Args) {
		return fmt.Errorf("Expected %d argument(s), but received %d",
			len(c.Args), argIndex)
	}
	return nil
}

// Option represents an optional parameter that can be passed as part of a
// Rocket command
type Option struct {
	Key      string
	HelpText string
	Format   *regexp.Regexp
	Value    string
}

// validate returns nil if the given value meets the format requirements for
// this option, returns the validation error otherwise.
func (o *Option) validate(value string) error {
	// Check that the value meets the required format
	if !o.Format.MatchString(value) {
		return fmt.Errorf("Invalid format for option \"%s\"."+
			"Format must match regular expression %s.", o.Key, o.Format.String())
	}
	return nil
}

// Argument represents a required parameter that Rocket will check as part of a command.
type Argument struct {
	Name     string
	HelpText string
	Format   *regexp.Regexp
	Value    string
}

// validate returns nil if the given value meets the format requirements for
// this option, returns the validation error otherwise.
func (a *Argument) validate(value string) error {
	// Check that the value meets the required format
	if !a.Format.MatchString(value) {
		return fmt.Errorf("Invalid format for argument \"%s\"."+
			"Format must match regular expression %s.", a.Name, a.Format.String())
	}
	return nil
}

// Context stores a Slack message and the user who sent it.
type Context struct {
	Message *slack.Msg
	User    model.Member
	Options map[string]Option
	Args    []Argument
}

// CommandHandler is the interface all handlers of Rocket commands must implement.
type CommandHandler func(Context) (string, slack.PostMessageParameters)
