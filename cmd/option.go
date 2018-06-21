package cmd

import (
	"fmt"
	"regexp"
)

// Option represents a parameter that can be passed as part of a
// Rocket command
type Option struct {
	// Key is this option's identifier. Of course, keys for different options
	// under the same command should always be unique. For exmaple, one might
	// create a command with one option who's key is `name`. In this case the
	// user would assign a value to this key in their Slack command with
	// `name={myvalue}`.
	Key   string
	Value string

	// HelpText is a description of what the option is used for.
	HelpText string

	// Format is a `regexp.Regexp` object that specifies the required format of
	// a value for an option. The `cmd` framework will enforce that this format
	// is met when a user enters a value for a given option, and will return an
	// appropriate error response if this is not the case. Commonly used format
	// `Regex`s can be found in [bot/util.go](bot/util.go).
	Format *regexp.Regexp

	// Required defines whether or not a value for this option is required when
	// a user uses this command. The `cmd` framework will enforce that a value
	// is set for each required option when a user enters a command, and will
	// return an appropriate error if this is not the case.
	Required bool
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
