package model

import "time"
import "github.com/nlopes/slack"

type Member struct {
	TableName struct{} `sql:"members" json:"-"`

	SlackID        string    `sql:",pk" json:"-"`
	Name           string    `json:"name"`
	Email          string    `json:"-"`
	GithubUsername string    `json:"githubUsername"`
	Major          string    `sql:"program" json:"major"`
	Position       string    `json:"position"`
	ImageURL       string    `json:"imageUrl"`
	IsAdmin        bool      `json:"-"`
	CreatedAt      time.Time `json:"-"`
}

type Members []*Member

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
	}

	if len(m.Name) == 0 {
		attachments[0].Color = "danger"
	}
	if len(m.Email) == 0 {
		attachments[1].Color = "danger"
	}
	if len(m.Position) == 0 {
		attachments[2].Color = "danger"
	}
	if len(m.GithubUsername) == 0 {
		attachments[3].Color = "danger"
	}
	if len(m.Major) == 0 {
		attachments[4].Color = "danger"
	}

	return attachments
}
