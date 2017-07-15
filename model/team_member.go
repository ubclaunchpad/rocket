package model

type TeamMember struct {
	TeamName      string `sql:",pk"`
	MemberSlackID string `sql:",pk"`

	Team   *Team   `sql:"-"`
	Member *Member `sql:"-"`
}
