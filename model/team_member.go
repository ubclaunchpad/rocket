package model

type TeamMember struct {
	TeamID   string `sql:",pk"`
	MemberID string `sql:",pk"`

	Team   *Team   `sql:"-"`
	Member *Member `sql:"-"`
}
