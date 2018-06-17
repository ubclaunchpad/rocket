package data

import (
	"errors"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/config"
)

// DAL represents the data abstraction layer and provides an interface to the
// database. This is just a wrapper around the PG database object.
type DAL struct {
	db orm.DB
}

// New returns a new DAL instance based on a configuration object.
func New(c *config.Config) *DAL {
	opts := &pg.Options{
		Addr:            c.PostgresHost + ":" + c.PostgresPort,
		User:            c.PostgresUser,
		Password:        c.PostgresPass,
		Database:        c.PostgresDatabase,
		MaxRetries:      10,
		MinRetryBackoff: time.Second,
		MaxRetryBackoff: time.Second * 10,
	}

	db := pg.Connect(opts)
	dal := &DAL{db}

	err := dal.Ping()
	if err != nil {
		log.WithError(err).Fatal("Error initializing the database")
	}

	return dal
}

// Ping checks that we can reach the database.
func (dal *DAL) Ping() error {
	i := 0
	_, err := dal.db.QueryOne(pg.Scan(&i), "SELECT 1")
	return err
}

// Close closes the connection to the database.
func (dal *DAL) Close() error {
	database, ok := dal.db.(*pg.DB)
	if !ok {
		return errors.New("dal.db not of type *pg.DB")
	}
	return database.Close()
}
