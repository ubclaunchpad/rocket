package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubclaunchpad/rocket/model"
)

func TestCreateGetAndDeleteTeam(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	dal, cleanupFunc, err := newTestDBConnection()
	assert.Nil(t, err)
	defer cleanupFunc()

	// Create a new team
	team := &model.Team{
		Name:         "team-bob",
		GithubTeamID: 1234,
		Platform:     "winning",
		Members: []*model.Member{&model.Member{
			SlackID:  "1234",
			Name:     "Little Bruno",
			Position: "A REAL GUY",
		}},
	}
	err = dal.CreateTeam(team)
	assert.Nil(t, err)

	// Get team by name
	teamGetByName := &model.Team{Name: "team-bob"}
	err = dal.GetTeamByName(teamGetByName)
	assert.Nil(t, err)
	assert.Equal(t, team.Platform, teamGetByName.Platform)

	// Get team by ID
	teamGetByID := &model.Team{GithubTeamID: 1234}
	err = dal.GetTeamByGithubID(teamGetByID)
	assert.Nil(t, err)
	assert.Equal(t, team.Platform, teamGetByID.Platform)

	// Delete team
	teamDeleteByID := &model.Team{Name: "team-bob"}
	err = dal.DeleteTeamByName(teamDeleteByID)
	assert.Nil(t, err)

	// Attempt to get
	err = dal.GetTeamByGithubID(teamGetByID)
	assert.NotNil(t, err)
}
