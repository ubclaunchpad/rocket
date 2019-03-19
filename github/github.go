package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ubclaunchpad/rocket/config"
	"golang.org/x/oauth2"

	gh "github.com/google/go-github/github"
)

// API provides a client to the GitHub API.
type API struct {
	organization string
	httpClient   *http.Client
	*gh.Client
	cache
}

type cache struct {
	validDuration time.Duration
	statsExpiry   time.Time
	statsData     *OrgStats
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
		cache{validDuration: time.Duration(6 * time.Hour)},
	}
}

// OrgStats represents basic stats about the configured organization's
// repositories and activity
type OrgStats struct {
	Repositories int            `json:"repositories"`
	Stargazers   int            `json:"stargazers"`
	Topics       map[string]int `json:"topics"`
	Languages    map[string]int `json:"languages"`

	CommitTotal int               `json:"commit_total"`
	CommitGraph map[time.Time]int `json:"commit_graph"`
}

// GetOrgStats collects basic stats about the configured organization's
// repositories and activity
func (api *API) GetOrgStats() (OrgStats, error) {
	// Return cache while valid
	if api.cache.statsExpiry.After(time.Now()) && api.cache.statsData != nil {
		return *api.cache.statsData, nil
	}

	// Generate new stats
	ctx := context.Background()
	repos, _, err := api.Repositories.ListByOrg(ctx, api.organization, nil)
	if err != nil {
		return OrgStats{}, err
	}

	// Collect stats from repositories
	stats := OrgStats{
		Topics:      make(map[string]int),
		Languages:   make(map[string]int),
		CommitGraph: make(map[time.Time]int),
	}
	for _, r := range repos {
		// Collect basic repository stats
		stats.Repositories++
		stats.Stargazers += r.GetStargazersCount()
		stats.Languages[r.GetLanguage()]++
		for _, t := range r.Topics {
			stats.Topics[t]++
		}

		// Collect activity stats
		activity, _, err := api.Repositories.ListCommitActivity(ctx, api.organization, r.GetName())
		if err != nil {
			continue
		}
		for _, week := range activity {
			stats.CommitGraph[week.GetWeek().Time] += week.GetTotal()
		}
	}

	// Store in cache
	api.cache.statsExpiry = time.Now().Add(api.cache.validDuration)
	api.cache.statsData = &stats
	return stats, nil
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

// RemoveTeam removes team with given ID from organization
func (api *API) RemoveTeam(id int) error {
	_, err := api.Organizations.DeleteTeam(context.Background(), id)
	return err
}
