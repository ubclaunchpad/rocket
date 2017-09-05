package github

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ubclaunchpad/rocket/config"
	"golang.org/x/oauth2"

	gh "github.com/jordanschalm/go-github/github"
)

type API struct {
	httpClient *http.Client
	*gh.Client
}

func New(c *config.Config) *API {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := gh.NewClient(tc)

	return &API{
		tc,
		client,
	}
}

// UserExists returns true if the user
func (api *API) UserExists(username string) (bool, error) {
	_, _, err := api.Users.Get(context.Background(), username)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (api *API) AddUserToTeam(username string, teamID int) error {
	_, _, err := api.Organizations.AddTeamMembership(
		context.Background(), teamID, username, nil,
	)
	return err
}

func (api *API) RemoveUserFromTeam(username string, teamID int) error {
	_, err := api.Organizations.RemoveTeamMembership(
		context.Background(), teamID, username,
	)
	return err
}

func (api *API) CreateTeam(name string) (*gh.Team, error) {
	log.Info("Getting teams")
	teams, _, err := api.Organizations.ListTeams(context.Background(), "ubclaunchpad", nil)
	if err != nil {
		return nil, err
	}

	log.Info("Checking if team exists")
	// Check if the team already exists
	for _, team := range teams {
		if *team.Name == name {
			return team, nil
		}
	}

	log.Info("Creating team")
	// Otherwise, create it
	team := &gh.Team{
		Name:    &name,
		Privacy: gh.String("closed"),
	}
	t, _, err := api.Organizations.CreateTeam(context.Background(), "ubclaunchpad", team)
	return t, err
}

func (api *API) RemoveTeam(id int) error {
	_, err := api.Organizations.DeleteTeam(context.Background(), id)
	return err
}
