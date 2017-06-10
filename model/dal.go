package model

import (
	"strconv"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/ubclaunchpad/rocket/config"

	"gopkg.in/pg.v4"
)

// DAL represents the data abstraction layer and provides an interface
// to the database.
type DAL struct {
	db pg.DB
}

var (
	instance DAL
	once     sync.Once
)

// Init initializes the DAL with a configuration. Init can only be called once.
func Init(c *config.Config) {
	once.Do(func() {
		db := pg.Connect(&pg.Options{
			Addr:     conf.PostgresHost + ":" + strconv.FormatUint(uint64(conf.PostgreSQLPort), 10),
			User:     conf.PostgresUsername,
			Password: conf.PostgresPassword,
			Database: conf.PostgresDatabase,
		})
		instance = DAL{db}
		err := instance.Ping()
		if err != nil {
			log.Fatal("Error initializing the database: ", err.Error())
		}
	})
}

// Ping checks that we can reach the database.
func (dal *DAL) Ping() error {
	i := 0
	_, err = dp.db.QueryOne(pg.Scan(&i), "SELECT 1")
	return err
}

// Close closes the connection to the database.
func (dal *DAL) Close() error {
	return dal.db.Close()
}
