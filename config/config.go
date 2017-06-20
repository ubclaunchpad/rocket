package config

import (
	"os"
)

// Config represents configuration option state for the app.
type Config struct {
	Token            string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPass     string
	PostgresDatabase string
	Host             string
	Port             string
}

// FromEnv creates a configuration from the environment.
func FromEnv() *Config {
	return &Config{
		Token:        os.Getenv("ROCKET_TOKEN"),
		PostgresHost: os.Getenv("ROCKET_POSTGRESHOST"),
		PostgresPort: os.Getenv("ROCKET_POSTGRESPORT"),
		PostgresUser: os.Getenv("ROCKET_POSTGRESUSER"),
		PostgresPass: os.Getenv("ROCKET_POSTGRESPASS"),
	}
}
