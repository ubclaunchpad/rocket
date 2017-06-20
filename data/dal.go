package data

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/ubclaunchpad/rocket/config"

	"github.com/go-pg/pg"
)

// DAL represents the data abstraction layer and provides an interface
// to the database.
type DAL struct {
	db *pg.DB
}

var (
	instance DAL
	once     sync.Once
)

// Init initializes the DAL with a configuration. Init can only be called once.
func Init(c *config.Config) {
	once.Do(func() {
		db := pg.Connect(&pg.Options{
			Addr:     c.PostgresHost + ":" + c.PostgresPort,
			User:     c.PostgresUser,
			Password: c.PostgresPass,
			Database: c.PostgresDatabase,
		})
		instance = DAL{db}
		err := instance.Ping()
		if err != nil {
			log.Fatal("Error initializing the database: ", err.Error())
		}
	})
}

func Get() *DAL {
	return &instance
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
