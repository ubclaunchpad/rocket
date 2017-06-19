package main

import (
	"fmt"
	"os"

	"github.com/nlopes/slack"
)

func main() {
	token := os.Getenv("ROCKET_TOKEN")
	// pgHost := os.Getenv("ROCKET_POSTGRESHOST")
	// pgPort := os.Getenv("ROCKET_POSTGRESPORT")
	// pgUsername := os.Getenv("ROCKET_POSTGRESUSERNAME")
	// pgPassword := os.Getenv("ROCKET_POSTGRESPASSWORD")

	api := slack.New(token)
	rtm := api.NewRTM()

	go rtm.ManageConnection()

	for evt := range rtm.IncomingEvents {
		switch evt.Data.(type) {
		case *slack.MessageEvent:
			msg := evt.Data.(*slack.MessageEvent).Msg
			fmt.Println("===")
			fmt.Println("Contents: ", msg.Text)
			fmt.Println("User / Username: ", msg.User, msg.Username)
			fmt.Println("Name / Members: ", msg.Name, msg.Members)
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello "+msg.Username+". I'm Rocket, your friendly neighbourhood Slack app. "+
				"I don't do much yet, but hopefully that will change soon :robot:", msg.Channel))
		}
	}
}
