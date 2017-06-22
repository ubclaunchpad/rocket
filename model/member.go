package model

import "time"

type Member struct {
	TableName struct{} `sql:"members" json:"-"`

	Email          string    `sql:",pk" json:"-"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	GithubUsername string    `json:"githubUsername"`
	Program        string    `json:"program"`
	ImageURL       string    `json:"imageUrl"`
	CreatedAt      time.Time `json:"-"`
}

type Members []*Member
