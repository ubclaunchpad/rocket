package data

import "github.com/ubclaunchpad/rocket/model"
import "github.com/go-pg/pg/orm"

func (dal *DAL) GetTeamByName(team *model.Team) error {
	return dal.db.Model(team).
		Where("name = ?name").
		Column("Members").
		Select()
}

func (dal *DAL) GetTeamByGithubID(team *model.Team) error {
	return dal.db.Model(team).
		Where("github_team_id = ?github_team_id").
		Column("Members").
		Select()
}

func (dal *DAL) GetTeams(teams *model.Teams) error {
	return dal.db.Model(teams).
		Column("Members").
		Relation("Members", func(q *orm.Query) (*orm.Query, error) {
			return q.Order("name ASC"), nil
		}).
		Order("name ASC").
		Select()
}

func (dal *DAL) CreateTeam(team *model.Team) error {
	_, err := dal.db.Model(team).
		OnConflict("DO NOTHING").
		Insert()
	return err
}

func (dal *DAL) DeleteTeamByName(team *model.Team) error {
	_, err := dal.db.Model(team).
		Where("name = ?name").
		Delete()
	return err
}
