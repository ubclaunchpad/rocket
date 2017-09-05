package github

import (
	"context"
	"net/http"

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
