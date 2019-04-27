package data

import (
	"fmt"

	"github.com/ubclaunchpad/rocket/model"
)

// CreateTeamMember inserts a team member into the database or returns an error.
func (dal *DAL) CreateTeamMember(member *model.TeamMember) error {
	_, err := dal.db.Model(member).
		OnConflict("DO NOTHING").
		Insert()
	return err
}

// GetTeamMember gets a team member by GithubUserId and MemberSlackID or returns
// an error.
func (dal *DAL) GetTeamMember(teamMember *model.TeamMember) error {
	return dal.db.Model(teamMember).
		Where("team_github_team_id = ?team_github_team_id").
		Where("member_slack_id = ?member_slack_id").
		Select()
}

// SetIsTechLead sets is_tech_lead to the value on the given team member or
// returns an error if there is no such team member.
func (dal *DAL) SetIsTechLead(teamMember *model.TeamMember) error {
	res, err := dal.db.Model(teamMember).
		Where("team_github_team_id = ?team_github_team_id AND " +
			"member_slack_id = ?member_slack_id").
		Set("is_tech_lead = ?is_tech_lead").
		Update()
	if res.RowsAffected() == 0 {
		return fmt.Errorf("No such team member member_slack_id=%s "+
			"team_github_team_id=%d",
			teamMember.MemberSlackID, teamMember.GithubTeamID)
	}
	return err
}

// GetTechLeads returns a list of all members who are tech leads, or returns
// an error.
func (dal *DAL) GetTechLeads() (*model.Members, error) {
	members := &model.Members{}
	err := dal.db.Model(members).
		Where("slack_id in (SELECT member_slack_id FROM team_members " +
			"WHERE is_tech_lead = true)").
		Order("name ASC").
		Select()
	return members, err
}

// GetTechLeadsByTeam returns a list of all members who are tech leads of the
// given team, or an error.
func (dal *DAL) GetTechLeadsByTeam(team *model.Team) (*model.Members, error) {
	members := &model.Members{}
	err := dal.db.Model(members).
		Where("slack_id in (SELECT member_slack_id FROM team_members "+
			"WHERE is_tech_lead = true and team_github_team_id = ?)",
			team.GithubTeamID).
		Order("name ASC").
		Select()
	return members, err
}

// DeleteTeamMember removes the given team member from database or
// returns an error.
func (dal *DAL) DeleteTeamMember(member *model.TeamMember) error {
	_, err := dal.db.Model(member).
		Where("team_github_team_id = ?team_github_team_id").
		Where("member_slack_id = ?member_slack_id").
		Delete()
	return err
}
