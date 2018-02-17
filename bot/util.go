package bot

import "regexp"

var (
	// anyRegex matches any non-empty string
	anyRegex = regexp.MustCompile(".+")
	// lowerAlphaRegex matches words containing only lower case letters
	lowerAlphaRegex = regexp.MustCompile("[a-z]")
	// alphaRegex matches words containing only letters
	alphaRegex = regexp.MustCompile("[a-zA-Z]")
	// nameRegex matches people's names
	nameRegex = regexp.MustCompile("^[a-zA-Z'-]+$")
	// emailRegex matches email addresses
	emailRegex = regexp.MustCompile("[a-zA-Z0-9._+]+@[a-zA-Z0-9._]+")
	// usernameRegex matches any Slack username
	usernameRegex = regexp.MustCompile("^[a-z0-9][a-z0-9._-]*$")
)

// Slack encodes user mentions slightly differently in the message objects
// that we receive from the RTM than they appear in the app. This function
// converts a plain username ID (a 12-character string) to a correctly
// formatted mention.
func toMention(username string) string {
	return "<@" + username + ">"
}

// Parses a mention and returns the ID of the user that was mentioned.
func parseMention(mention string) string {
	if len(mention) != 12 {
		return ""
	}
	return mention[2:11]
}
