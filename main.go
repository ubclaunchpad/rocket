package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/bot"
	"github.com/ubclaunchpad/rocket/config"
	"github.com/ubclaunchpad/rocket/server"

	"github.com/ubclaunchpad/rocket/data"
)

func main() {
	cfg := config.FromEnv()

	dal := data.New(cfg)

	srv := server.New(cfg, dal, log.WithField("service", "server"))
	go srv.Start()

	slack := bot.New(cfg, dal, log.WithField("service", "slack"))
	slack.Start()
}
