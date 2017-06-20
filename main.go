package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/nlopes/slack"
)

func main() {
	token := os.Getenv("ROCKET_TOKEN")
	// pgHost := os.Getenv("ROCKET_POSTGRESHOST")
	// pgPort := os.Getenv("ROCKET_POSTGRESPORT")
	// pgUsername := os.Getenv("ROCKET_POSTGRESUSERNAME")
	// pgPassword := os.Getenv("ROCKET_POSTGRESPASSWORD")

	// host = os.Getenv("ROCKET_HOST")
	// port = os.Getenv("ROCKET_PORT")

	api := slack.New(token)
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
				"Text": msg.Text,
				"User": msg.User,
			}).Info("Message")
			rtm.SendMessage(rtm.NewOutgoingMessage("Hi, I'm Rocket, your friendly neighbourhood Slack app. "+
				"I don't do much yet, but hopefully that will change soon :robot_face:", msg.Channel))
		}
	}
}
