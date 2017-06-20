package model

import "time"

type Team struct {
	ID        string    `sql:",pk" json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`

	Members []*Member `sql:"-" pg:",many2many:team_members,joinFK:Member"`
}

type Teams []*Team
