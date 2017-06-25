package bot

import "regexp"

var (
	emailRegex = regexp.MustCompile("[a-zA-Z0-9._+]+@[a-zA-Z0-9._]+")
)

func toMention(username string) string {
	return "<@" + username + ">"
}

func parseEmail(email string) string {
	return emailRegex.FindString(email)
}
