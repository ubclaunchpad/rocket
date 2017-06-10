package main

import (
	"fmt"
	"os"

	"github.com/nlopes/slack"
)

func main() {
	token := os.Getenv("ROCKET_TOKEN")
	pgHost := os.Getenv("ROCKET_POSTGRESHOST")
	pgPort := os.Getenv("ROCKET_POSTGRESPORT")
	pgUsername := os.Getenv("ROCKET_POSTGRESUSERNAME")
	pgPassword := os.Getenv("ROCKET_POSTGRESPASSWORD")

	api := slack.New(token)
	rtm := api.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch msg.Data.(type) {
		case *slack.MessageEvent:
			message := msg.Data.(*slack.MessageEvent).Msg
			fmt.Println(message.Text)
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello "+message.Channel, message.Channel))
		}
	}
}
