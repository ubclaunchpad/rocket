package model

import "time"

// Member represents a team member
type Member struct {
	FirstName string
	LastName  string
	Email     string
	Program   string
	CreatedAt time.Time
}
