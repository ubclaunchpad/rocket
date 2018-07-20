package cmd

import "regexp"

var (
	// AnyRegex matches any non-empty string
	AnyRegex = regexp.MustCompile(".+")
	// AlphaRegex matches words containing only letters
	AlphaRegex = regexp.MustCompile("[a-zA-Z]")
	// EmailRegex matches email addresses
	EmailRegex = regexp.MustCompile("[a-zA-Z0-9._+]+@[a-zA-Z0-9._]+")
)

// ToMention converts a Slack username to a mention.
// Slack encodes user mentions slightly differently in the message objects
// that we receive from the RTM than they appear in the app. This function
// converts a plain username ID (a 12-character string) to a correctly
// formatted mention.
func ToMention(username string) string {
	return "<@" + username + ">"
}

// ParseMention parses a mention and returns the ID of the user that was mentioned.
func ParseMention(mention string) string {
	if len(mention) != 12 {
		return ""
	}
	return mention[2:11]
}
