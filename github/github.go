package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ubclaunchpad/rocket/config"
	"golang.org/x/oauth2"

	gh "github.com/google/go-github/github"
)

// API provides a client to the GitHub API.
type API struct {
	organization string
	httpClient   *http.Client
	*gh.Client
}

// New creates and returns a GitHub API object based on a configuration object,
// configured for use with the given organization
func New(organization string, c *config.Config) *API {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := gh.NewClient(tc)

	return &API{
		organization,
		tc,
		client,
	}
}

// UserExists checks if a given user exists in Github
func (api *API) UserExists(username string) (bool, error) {
	_, _, err := api.Users.Get(context.Background(), username)
	if err != nil {
		return false, err
	}
	return true, nil
}

// AddUserToTeam adds given user to given team
func (api *API) AddUserToTeam(username string, teamID int) error {
	_, _, err := api.Organizations.AddTeamMembership(
		context.Background(), teamID, username, nil,
	)
	return err
}

// RemoveUserFromOrg removes given user from configured organization
func (api *API) RemoveUserFromOrg(username string) error {
	_, err := api.Organizations.RemoveOrgMembership(
		context.Background(), username, api.organization,
	)
	return err
}

// RemoveUserFromTeam removes given user from configured organization
func (api *API) RemoveUserFromTeam(username string, teamID int) error {
	_, err := api.Organizations.RemoveTeamMembership(
		context.Background(), teamID, username,
	)
	return err
}

// CreateTeam creates a team in the configured organization
func (api *API) CreateTeam(name string) (*gh.Team, error) {
	teams, _, err := api.Organizations.ListTeams(context.Background(), api.organization, nil)
	if err != nil {
		return nil, err
	}

	// Check if the team already exists
	for _, team := range teams {
		if *team.Name == name {
			return team, nil
		}
	}

	// Otherwise, create it
	team := &gh.NewTeam{
		Name:    name,
		Privacy: gh.String("closed"),
	}
	t, _, err := api.Organizations.CreateTeam(context.Background(), api.organization, team)
	return t, err
}

// GetTeam retrieves team with given team ID from configured organization
func (api *API) GetTeam(id int) (*gh.Team, error) {
	teams, _, err := api.Organizations.ListTeams(context.Background(), api.organization, nil)
	if err != nil {
		return nil, err
	}
	for _, team := range teams {
		if *team.ID == id {
			return team, nil
		}
	}
	return nil, fmt.Errorf("GitHub team with ID %d not found", id)
}

func (api *API) RemoveTeam(id int) error {
	_, err := api.Organizations.DeleteTeam(context.Background(), id)
	return err
}
