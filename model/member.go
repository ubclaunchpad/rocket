package model

import (
	"time"

	"github.com/nlopes/slack"
)

// Member is the concrete representation of a Launch Pad club member in the
// database.
type Member struct {
	TableName struct{} `sql:"members" json:"-"`

	SlackID        string    `sql:",pk" json:"-"`
	Name           string    `json:"name"`
	Email          string    `json:"-"`
	GithubUsername string    `json:"githubUsername"`
	Major          string    `sql:"program" json:"major"`
	Position       string    `json:"position"`
	Biography      string    `json:"biography"`
	ImageURL       string    `json:"imageUrl"`
	IsTechLead     bool      `json:"isTechLead"`
	IsAdmin        bool      `json:"-"`
	CreatedAt      time.Time `json:"-"`
}

// Members is a list of members
type Members []*Member

// SlackAttachments creates and returns a set of Slack attachments (strictly
// for use in messages sent to Slack clients) that describe the member's
// profile. Each profile field is one attachment, and is colour-coded based
// on whether it's been filled yet.
func (m *Member) SlackAttachments() []slack.Attachment {
	attachments := []slack.Attachment{
		slack.Attachment{
			Text:  "Name: " + m.Name,
			Color: "good",
		},
		slack.Attachment{
			Text:  "Email: " + m.Email,
			Color: "good",
		},
		slack.Attachment{
			Text:  "Position: " + m.Position,
			Color: "good",
		},
		slack.Attachment{
			Text:  "GitHub Username: " + m.GithubUsername,
			Color: "good",
		},
		slack.Attachment{
			Text:  "Major: " + m.Major,
			Color: "good",
		},
		slack.Attachment{
			Text:  "Biography: " + m.Biography,
			Color: "good",
		},
	}

	for _, attachment := range attachments {
		if len(attachment.Text) == 0 {
			attachment.Color = "danger"
		}
	}

	if m.IsTechLead {
		attachments = append(attachments, slack.Attachment{
			Text:  "Tech Lead",
			Color: "good",
		})
	}

	return attachments
}
