package data

import "github.com/ubclaunchpad/rocket/model"

func (dal *DAL) GetTeamById(team *model.Team) error {
	return dal.db.Model(team).
		Where("id = ?id").
		Column("Members").
		Select()
}

func (dal *DAL) GetTeams(teams *model.Teams) error {
	return dal.db.Model(teams).
		Column("Members").
		Select()
}
