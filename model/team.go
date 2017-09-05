package model

import (
	"strings"
	"time"

	"github.com/nlopes/slack"
)

type Team struct {
	TableName struct{} `sql:"teams" json:"-"`

	Name         string    `json:"name"`
	GithubTeamID int       `sql:",pk" json:"-" pg:"github_team_id"`
	CreatedAt    time.Time `json:"-"`

	Members []*Member `sql:"-" json:"members" pg:",many2many:team_members,joinFK:Member"`
}

type Teams []*Team

func (t *Team) SlackAttachments() []slack.Attachment {
	members := []string{}
	for _, member := range t.Members {
		members = append(members, member.Name)
	}
	membersString := strings.Join(members, ", ")

	attachments := []slack.Attachment{
		slack.Attachment{
			Text:  "Name: " + t.Name,
			Color: "good",
		},
		slack.Attachment{
			Text:  "Members: " + membersString,
			Color: "good",
		},
	}

	return attachments
}
