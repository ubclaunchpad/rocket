package bot

import "regexp"

var (
	emailRegex = regexp.MustCompile("[a-zA-Z0-9._+]+@[a-zA-Z0-9._]+")
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

// Slack also encodes emails slightly differently. This function parses an
// actual email from an Slack-formatted email.
func parseEmail(email string) string {
	return emailRegex.FindString(email)
}
