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

	// channels, err := api.GetChannels(true)
	// if err != nil {
	// 	log.WithError(err).Error("Could not view channels")
	// } else {
	// 	for _, c := range channels {
	// 		log.WithFields(log.Fields{
	// 			"ID":   c.ID,
	// 			"Name": c.Name,
	// 		}).Info()
	// 	}
	// }

	// ims, err := api.GetIMChannels()
	// if err != nil {
	// 	log.WithError(err).Error("Could not view IM channels")
	// } else {
	// 	for _, im := range ims {
	// 		log.WithFields(log.Fields{
	// 			"ID":   im.ID,
	// 			"User": im.User,
	// 		})
	// 	}
	// }

	// for evt := range rtm.IncomingEvents {
	// 	switch evt.Data.(type) {
	// 	case *slack.MessageEvent:
	// 		msg := evt.Data.(*slack.MessageEvent).Msg
	// 		log.WithFields(log.Fields{
	// 			"Text":    msg.Text,
	// 			"User":    msg.User,
	// 			"Channel": msg.Channel,
	// 			"Type":    msg.Type,
	// 		}).Info("Message")

	// 		tokens := strings.Split(msg.Text, " ")
	// 		if tokens[0] == mention(botUsername) {
	// 			member := model.Member{
	// 				SlackID: msg.User,
	// 			}
	// 			if tokens[1] == "me" {
	// 				if err := dal.GetMemberBySlackID(&member); err != nil {
	// 					log.WithError(err).Error("Error retrieving member by Slack ID")
	// 				} else {
	// 					params := slack.PostMessageParameters{
	// 						Attachments: member.SlackAttachments(),
	// 					}
	// 					api.PostMessage(msg.Channel, "Your UBC Launch Pad profile", params)
	// 				}
	// 			}
	// 			if tokens[1] == "init" {
	// 				member := model.Member{
	// 					SlackID: msg.User,
	// 				}
	// 				dal.CreateMember(&member)
	// 				api.PostMessage(msg.Channel, "I've set up your profile! Please use these commands to add information:\n"+
	// 					"`@rocket set name`\n`@rocket set email`\n`@rocket set github`\n`@rocket set program`", slack.PostMessageParameters{})
	// 			}
	// 			if tokens[1] == "set" {
	// 				if len(tokens) < 4 {
	// 					break
	// 				}
	// 				if tokens[2] == "name" {
	// 					member.Name = strings.Join(tokens[3:], " ")
	// 					if err := dal.SetMemberName(&member); err != nil {
	// 						api.PostMessage(msg.Channel, "An error occurred :cry:", slack.PostMessageParameters{})
	// 					} else {
	// 						api.PostMessage(msg.Channel, "Your name has been updated! :simple_smile:", slack.PostMessageParameters{})
	// 					}
	// 				}
	// 				if tokens[2] == "email" {
	// 					member.Email = tokens[3]
	// 					if err := dal.SetMemberEmail(&member); err != nil {
	// 						api.PostMessage(msg.Channel, "An error occurred :cry:", slack.PostMessageParameters{})
	// 					} else {
	// 						api.PostMessage(msg.Channel, "Your email has been updated! :simple_smile:", slack.PostMessageParameters{})
	// 					}
	// 				}
	// 				if tokens[2] == "github" {
	// 					member.GithubUsername = tokens[3]
	// 					if err := dal.SetMemberGitHubUsername(&member); err != nil {
	// 						api.PostMessage(msg.Channel, "An error occurred :cry:", slack.PostMessageParameters{})
	// 					} else {
	// 						api.PostMessage(msg.Channel, "Your GitHub username has been updated! :simple_smile:", slack.PostMessageParameters{})
	// 					}
	// 				}
	// 				if tokens[2] == "program" {
	// 					member.Program = tokens[3]
	// 					if err := dal.SetMemberProgram(&member); err != nil {
	// 						api.PostMessage(msg.Channel, "An error occurred :cry:", slack.PostMessageParameters{})
	// 					} else {
	// 						api.PostMessage(msg.Channel, "Your program has been updated! :simple_smile:", slack.PostMessageParameters{})
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}
