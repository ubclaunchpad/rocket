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
