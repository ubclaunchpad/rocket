package data

import (
	"errors"
	"testing"

	"github.com/go-pg/pg"
	"github.com/stretchr/testify/assert"
	"github.com/ubclaunchpad/rocket/config"
)

func newTestDBConnection() (*DAL, func(), error) {
	// Connect to database
	dal := New(&config.Config{
		PostgresHost:     "localhost",
		PostgresPort:     "5433",
		PostgresDatabase: "rocket_test_db",
		PostgresUser:     "rocket_test",
		PostgresPass:     "rickroll",
	})

	// Begin transaction
	database, ok := dal.db.(*pg.DB)
	if !ok {
		return nil, func() {}, errors.New("dal.db is not of type *pg.DB")
	}
	tx, err := database.Begin()
	if err != nil {
		return nil, func() {}, err
	}
	dal.db = tx

	// Return DAL with callback to roll back everything that happened in test
	return dal, func() {
		tx.Rollback()
		database.Close()
	}, nil
}

func TestNewDAL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	dal, cleanupFunc, err := newTestDBConnection()
	assert.Nil(t, err)
	defer cleanupFunc()

	err = dal.Ping()
	assert.Nil(t, err)
}
