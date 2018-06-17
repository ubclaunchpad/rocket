package data

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubclaunchpad/rocket/config"
)

func newTestDBConnection() *DAL {
	cfg := &config.Config{
		PostgresHost:     "localhost",
		PostgresPort:     "5432",
		PostgresDatabase: "rocket_test_db",
	}
	if os.Getenv("TRAVIS") == "true" {
		cfg.PostgresUser = "postgres"
	} else {
		cfg.PostgresUser = "rocket_test"
	}
	return New(cfg)
}

func TestNewDAL(t *testing.T) {
	dal := newTestDBConnection()
	defer dal.Close()
	err := dal.Ping()
	assert.Nil(t, err)
}
