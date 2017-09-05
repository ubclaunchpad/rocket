package model

type TeamMember struct {
	GithubTeamID  int    `sql:",pk" pg:"team_github_team_id"`
	MemberSlackID string `sql:",pk"`

	Team   *Team   `sql:"-"`
	Member *Member `sql:"-"`
}
