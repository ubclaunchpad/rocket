package data

import (
	"github.com/go-pg/pg/orm"
	"github.com/ubclaunchpad/rocket/model"
)

// GetTeamByName provides team with corresponding name
func (dal *DAL) GetTeamByName(team *model.Team) error {
	return dal.db.Model(team).
		Where("name = ?name").
		Column("Members").
		Select()
}

// GetTeamByGithubID provides team with corresponding GitHub ID
func (dal *DAL) GetTeamByGithubID(team *model.Team) error {
	return dal.db.Model(team).
		Where("github_team_id = ?github_team_id").
		Column("Members").
		Select()
}

// GetTeams gets all current teams
func (dal *DAL) GetTeams(teams *model.Teams) error {
	return dal.db.Model(teams).
		Column("Members").
		Relation("Members", func(q *orm.Query) (*orm.Query, error) {
			return q.Order("name ASC"), nil
		}).
		Order("name ASC").
		Select()
}

// GetTeamNames gets all names of current teams
func (dal *DAL) GetTeamNames(teams *model.Teams) error {
	return dal.db.Model(teams).
		Column("name").
		Select()
}

// CreateTeam inserts given team into the database
func (dal *DAL) CreateTeam(team *model.Team) error {
	_, err := dal.db.Model(team).
		OnConflict("DO NOTHING").
		Insert()
	return err
}

// UpdateTeam updates given team with new team
func (dal *DAL) UpdateTeam(currentTeam, newTeam *model.Team) error {
	// Only update values if they were set
	if newTeam.Name != "" {
		currentTeam.Name = newTeam.Name
	}
	if newTeam.Platform != "" {
		currentTeam.Platform = newTeam.Platform
	}
	_, err := dal.db.Model(currentTeam).
		Update(
			"name", currentTeam.Name,
			"platform", currentTeam.Platform)
	return err
}

// DeleteTeamByName deletes team with given name from the database
func (dal *DAL) DeleteTeamByName(team *model.Team) error {
	_, err := dal.db.Model(team).
		Where("name = ?name").
		Delete()
	return err
}
