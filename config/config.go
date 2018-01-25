package config

import (
	"os"
)

// Config represents configuration options for the app.
type Config struct {
	Host             string
	Port             string
	SlackToken       string
	GithubToken      string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPass     string
	PostgresDatabase string
}

// FromEnv creates and returns a configuration object from the environment.
func FromEnv() *Config {
	return &Config{
		Host:         os.Getenv("ROCKET_HOST"),
		Port:         os.Getenv("ROCKET_PORT"),
		SlackToken:   os.Getenv("ROCKET_SLACKTOKEN"),
		GithubToken:  os.Getenv("ROCKET_GITHUBTOKEN"),
		PostgresHost: os.Getenv("ROCKET_POSTGRESHOST"),
		PostgresPort: os.Getenv("ROCKET_POSTGRESPORT"),
		PostgresUser: os.Getenv("ROCKET_POSTGRESUSER"),
		PostgresPass: os.Getenv("ROCKET_POSTGRESPASS"),
	}
}
