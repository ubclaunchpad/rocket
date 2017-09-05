package data

import "github.com/ubclaunchpad/rocket/model"

func (dal *DAL) GetTeamByName(team *model.Team) error {
	return dal.db.Model(team).
		Where("name = ?name").
		Column("Members").
		Select()
}

func (dal *DAL) GetTeams(teams *model.Teams) error {
	return dal.db.Model(teams).
		Column("Members").
		Select()
}

func (dal *DAL) CreateTeam(team *model.Team) error {
	_, err := dal.db.Model(team).
		OnConflict("DO NOTHING").
		Insert()
	return err
}

func (dal *DAL) SetTeamGithubTeamName(team *model.Team) error {
	_, err := dal.db.Model(team).
		Set("github_team_name = ?github_team_name").
		Update()
	return err
}

func (dal *DAL) SetTeamGithubTeamID(team *model.Team) error {
	_, err := dal.db.Model(team).
		Set("github_team_id = ?github_team_id").
		Update()
	return err
}

func (dal *DAL) DeleteTeam(team *model.Team) error {
	_, err := dal.db.Model(team).
		Where("name = ?name").
		Delete()
	return err
}
