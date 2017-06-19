package model

import "time"

type Member struct {
	ID        string `sql:",pk"`
	Email     string
	FirstName string
	LastName  string
	Program   string
	ImageURL  string
	CreatedAt time.Time

	Teams []*Team `sql:"-" pg:",many2many:team_members,joinFK:Team"`
}
