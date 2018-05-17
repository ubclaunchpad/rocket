package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/bot"
	"github.com/ubclaunchpad/rocket/config"
	"github.com/ubclaunchpad/rocket/github"
	"github.com/ubclaunchpad/rocket/plugin"
	"github.com/ubclaunchpad/rocket/server"

	"github.com/ubclaunchpad/rocket/data"
)

func main() {
	// Create a configuration object by pulling the value of a bunch of
	// environment variables.
	cfg := config.FromEnv()

	// Connect the database and initialize the data access layer. We use the
	// URL, database, and password specified in the config. This will panic
	// if we fail to connect to the database.
	dal := data.New(cfg)

	// Create a client to the GitHub API, using the token from the config.
	gh := github.New(cfg)

	// Set up a server listening on the interface specified in the
	// config. This will panic if the server fails to bind to the interface
	// or dies for any reason after beginning listening.
	srv := server.New(cfg, dal, log.WithField("service", "server"))

	// Set up the Slack bot. This will create an RTM that receives
	// events from Slack and respond to them as needed.
	slackBot := bot.New(cfg, dal, gh, log.WithField("service", "slack"))

	// Load plugins
	plugin.RegisterPlugins(slackBot)

	// Start Slack bot and HTTP server
	go srv.Start()
	slackBot.Start()
}
