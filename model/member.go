package model

import "time"
import "github.com/nlopes/slack"

type Member struct {
	TableName struct{} `sql:"members" json:"-"`

	SlackID        string    `sql:",pk" json:"-"`
	Name           string    `json:"name"`
	Email          string    `json:"-"`
	GithubUsername string    `json:"githubUsername"`
	Program        string    `json:"program"`
	ImageURL       string    `json:"imageUrl"`
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
			Text:  "GitHub Username: " + m.GithubUsername,
			Color: "good",
		},
		slack.Attachment{
			Text:  "Program: " + m.Program,
			Color: "good",
		},
	}

	if len(m.Name) == 0 {
		attachments[0].Color = "danger"
	}
	if len(m.Email) == 0 {
		attachments[1].Color = "danger"
	}
	if len(m.GithubUsername) == 0 {
		attachments[2].Color = "danger"
	}
	if len(m.Program) == 0 {
		attachments[3].Color = "danger"
	}

	return attachments
}
