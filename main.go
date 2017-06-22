package main

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/config"
	"github.com/ubclaunchpad/rocket/server"

	"github.com/nlopes/slack"
	"github.com/ubclaunchpad/rocket/data"
)

func main() {
	cfg := config.FromEnv()

	data.Init(cfg)
	dal := data.Get()

	srv := server.New(cfg, dal)
	go srv.Start()

	api := slack.New(cfg.SlackToken)
	rtm := api.NewRTM()

	go rtm.ManageConnection()

	channels, err := api.GetChannels(true)
	if err != nil {
		log.WithError(err).Error("Could not view channels")
	} else {
		for _, c := range channels {
			log.WithFields(log.Fields{
				"ID":   c.ID,
				"Name": c.Name,
			}).Info()
		}
	}

	for evt := range rtm.IncomingEvents {
		switch evt.Data.(type) {
		case *slack.MessageEvent:
			msg := evt.Data.(*slack.MessageEvent).Msg
			log.WithFields(log.Fields{
				"Text":    msg.Text,
				"User":    msg.User,
				"Channel": msg.Channel,
				"Type":    msg.Type,
			}).Info("Message")
			if strings.Index(msg.Text, "rocket") >= 0 {
				rtm.SendMessage(rtm.NewOutgoingMessage("Hi, I'm Rocket, _your_ friendly neighbourhood Slack app. "+
					"I don't do much yet, but hopefully that will change soon :robot_face:", msg.Channel))
				api.PostMessage(msg.Channel, "Hello _there_, *what* is happening, `code`\n```\nlots of\ncode\n```", slack.PostMessageParameters{})
			}
		}
	}
}
