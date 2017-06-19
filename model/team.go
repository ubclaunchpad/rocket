package model

import "time"

type Team struct {
	ID        string `sql:",pk"`
	Name      string
	CreatedAt time.Time

	Members []*Member `sql:"-" pg:",many2many:team_members,joinFK:Member"`
}
