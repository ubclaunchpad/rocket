package data

import "github.com/ubclaunchpad/rocket/model"

func (dal *DAL) CreateTeamMember(member *model.TeamMember) error {
	_, err := dal.db.Model(member).
		OnConflict("DO NOTHING").
		Insert()
	return err
}

func (dal *DAL) DeleteTeamMember(member *model.TeamMember) error {
	_, err := dal.db.Model(member).
		Where("team_name = ?team_name").
		Where("member_slack_id = ?member_slack_id").
		Delete()
	return err
}
