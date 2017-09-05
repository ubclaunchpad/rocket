package config

import (
	"os"
)

// Config represents configuration option state for the app.
type Config struct {
	Host             string
	Port             string
	Domain           string
	SlackToken       string
	GithubToken      string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPass     string
	PostgresDatabase string
}

// FromEnv creates a configuration from the environment.
func FromEnv() *Config {
	return &Config{
		Host:         os.Getenv("ROCKET_HOST"),
		Port:         os.Getenv("ROCKET_PORT"),
		Domain:       os.Getenv("ROCKET_DOMAIN"),
		SlackToken:   os.Getenv("ROCKET_SLACKTOKEN"),
		GithubToken:  os.Getenv("ROCKET_GITHUBTOKEN"),
		PostgresHost: os.Getenv("ROCKET_POSTGRESHOST"),
		PostgresPort: os.Getenv("ROCKET_POSTGRESPORT"),
		PostgresUser: os.Getenv("ROCKET_POSTGRESUSER"),
		PostgresPass: os.Getenv("ROCKET_POSTGRESPASS"),
	}
}
