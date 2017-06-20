package model

import "time"

type Member struct {
	ID             string    `sql:",pk" json:"-"`
	Email          string    `json:"email"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	GithubUsername string    `json:"githubUsername"`
	Program        string    `json:"program"`
	ImageURL       string    `json:"imageUrl"`
	CreatedAt      time.Time `json:"-"`

	Teams []*Team `sql:"-" pg:",many2many:team_members,joinFK:Team"`
}

type Members []*Member
