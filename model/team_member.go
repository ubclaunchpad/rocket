package model

// TeamMember represents the concrete relationship between teams and members
// in the database.
type TeamMember struct {
	GithubTeamID  int    `sql:"team_github_team_id,pk" json:"-"`
	MemberSlackID string `sql:",pk" json:"-"`
	IsTechLead    bool   `json:"isTechLead"`

	Team   *Team   `sql:"-"`
	Member *Member `sql:"-"`
}
