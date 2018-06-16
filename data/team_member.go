package data

import "github.com/ubclaunchpad/rocket/model"

// CreateTeamMember inserts a team member into the database
func (dal *DAL) CreateTeamMember(member *model.TeamMember) error {
	_, err := dal.db.Model(member).
		OnConflict("DO NOTHING").
		Insert()
	return err
}

// DeleteTeamMember removes team member from database
func (dal *DAL) DeleteTeamMember(member *model.TeamMember) error {
	_, err := dal.db.Model(member).
		Where("team_github_team_id = ?team_github_team_id").
		Where("member_slack_id = ?member_slack_id").
		Delete()
	return err
}
