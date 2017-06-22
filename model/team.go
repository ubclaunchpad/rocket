package model

import "time"

type Team struct {
	TableName struct{} `sql:"teams" json:"-"`

	Name      string    `sql:",pk" json:"name"`
	CreatedAt time.Time `json:"-"`

	Members []*Member `sql:"-" pg:",many2many:team_members,joinFK:Member"`
}

type Teams []*Team
