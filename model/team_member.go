package model

// TeamMember represents the concrete relationship between teams and mbmers
// in the database.
type TeamMember struct {
	GithubTeamID  int    `sql:"team_github_team_id,pk"`
	MemberSlackID string `sql:",pk"`

	Team   *Team   `sql:"-"`
	Member *Member `sql:"-"`
}
