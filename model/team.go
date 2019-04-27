package model

import (
	"strings"
	"time"

	"github.com/nlopes/slack"
)

// Team represents the concrete representation of a team in the database.
type Team struct {
	TableName struct{} `sql:"teams" json:"-"`

	Name         string    `json:"name"`
	GithubTeamID int       `sql:",pk" json:"-" pg:"github_team_id"`
	Platform     string    `json:"platform" pg:"platform"`
	CreatedAt    time.Time `json:"-"`

	Members []*Member `sql:"-" json:"members" pg:",many2many:team_members,joinFK:Member"`
}

// Teams is a list of teams
type Teams []*Team

// SlackAttachments creates and returns a set of Slack attachments (strictly
// for use in messages sent to Slack clients) that describe the team's name
// and list of members.
func (t *Team) SlackAttachments(techLeads Members) []slack.Attachment {
	members := []string{}
	leads := []string{}
	for _, member := range t.Members {
		members = append(members, member.Name)
	}
	for _, lead := range techLeads {
		leads = append(leads, lead.Name)
	}
	membersString := strings.Join(members, ", ")
	leadsString := strings.Join(leads, ", ")

	attachments := []slack.Attachment{
		slack.Attachment{
			Text:  "Name: " + t.Name,
			Color: "good",
		},
		slack.Attachment{
			Text:  "Platform: " + t.Platform,
			Color: "good",
		},
		slack.Attachment{
			Text:  "Leads: " + leadsString,
			Color: "good",
		},
		slack.Attachment{
			Text:  "Members: " + membersString,
			Color: "good",
		},
	}

	return attachments
}
