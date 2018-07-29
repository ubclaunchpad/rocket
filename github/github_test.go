package github

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	gh "github.com/google/go-github/github"
)

func TestAPI_GetOrgStats(t *testing.T) {
	mockOrgStats := OrgStats{
		Repositories: 100,
		Topics:       make(map[string]int),
		Languages:    make(map[string]int),
		CommitGraph:  make(map[time.Time]int),
	}
	type fields struct {
		organization string
		httpClient   *http.Client
		Client       *gh.Client
		cache        cache
	}
	tests := []struct {
		name    string
		fields  fields
		want    OrgStats
		wantErr bool
	}{
		{
			"should get from cache if not expired",
			fields{"ubclaunchpad", nil, nil, cache{
				statsExpiry: time.Now().Add(time.Duration(6 * time.Hour)),
				statsData:   &mockOrgStats}},
			mockOrgStats,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &API{
				organization: tt.fields.organization,
				httpClient:   tt.fields.httpClient,
				Client:       tt.fields.Client,
				cache:        tt.fields.cache,
			}
			got, err := api.GetOrgStats()
			if (err != nil) != tt.wantErr {
				t.Errorf("API.GetOrgStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("API.GetOrgStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
