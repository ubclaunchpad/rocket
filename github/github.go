package github

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ubclaunchpad/rocket/config"
	"golang.org/x/oauth2"

	gh "github.com/google/go-github/github"
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

func (api *API) CreateTeam(name string) (*gh.Team, error) {
	teams, _, err := api.Organizations.ListTeams(context.Background(), "ubclaunchpad", nil)
	log.Info("listteams ", teams, err)
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
	team := &gh.Team{
		Name: &name,
	}
	t, _, err := api.Organizations.CreateTeam(context.Background(), "ubclaunchpad", team)
	log.Info("crateteam ", t, err)
	return t, err
}
