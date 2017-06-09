package main

import (
	"fmt"
	"os"

	"github.com/nlopes/slack"
)

func main() {
	api := slack.New(os.Getenv("ROCKET_TOKEN"))
	rtm := api.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch msg.Data.(type) {
		case *slack.MessageEvent:
			message := msg.Data.(*slack.MessageEvent).Msg
			fmt.Println(message.Text)
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello #general", "#general"))
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello "+message.Channel, message.Channel))
		}
	}
}
