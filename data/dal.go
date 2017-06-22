package data

import (
	log "github.com/sirupsen/logrus"

	"github.com/ubclaunchpad/rocket/config"

	"github.com/go-pg/pg"
)

// DAL represents the data abstraction layer and provides an interface
// to the database.
type DAL struct {
	db *pg.DB
}

// New returns a new DAL instance based on the config.
func New(c *config.Config) *DAL {
	db := pg.Connect(&pg.Options{
		Addr:     c.PostgresHost + ":" + c.PostgresPort,
		User:     c.PostgresUser,
		Password: c.PostgresPass,
		Database: c.PostgresDatabase,
	})
	dal := &DAL{db}

	err := dal.Ping()
	if err != nil {
		log.Fatal("Error initializing the database: ", err.Error())
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
	return dal.db.Close()
}
