package bot

import (
	"strings"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/config"
	"github.com/ubclaunchpad/rocket/data"
	"github.com/ubclaunchpad/rocket/model"
)

const (
	// Our Slack Bot's username
	username = "U5RU9TB38"

	// Message explaining commands
	helpMessage = "Hi there, I'm Rocket, Launch Pad's friendly neighbourhood Slack bot! :rocket:\n" +
		"You can create your profile with `@rocket init` and view your profile with `@rocket me`.\n" +
		"You can update your profile too!\n" +
		"```\n@rocket set name <name>\n@rocket set email\n@rocket set github\n@rocket set program\n```"

	// Message for when errors occur
	errorMessage = "Oops, an error occurred :robot_face:. Bruno must have coded a bug... Sorry about that!"
)

var (
	noParams = slack.PostMessageParameters{}
)

type Bot struct {
	token string
	api   *slack.Client
	rtm   *slack.RTM
	dal   *data.DAL
	log   *log.Entry
}

func New(cfg *config.Config, dal *data.DAL, log *log.Entry) *Bot {
	api := slack.New(cfg.SlackToken)
	return &Bot{
		token: cfg.SlackToken,
		api:   api,
		rtm:   api.NewRTM(),
		dal:   dal,
		log:   log,
	}
}

func (b *Bot) Start() {
	go b.rtm.ManageConnection()

	for evt := range b.rtm.IncomingEvents {
		switch evt.Data.(type) {
		case *slack.MessageEvent:
			b.handleMessageEvent(evt.Data.(*slack.MessageEvent).Msg)
		}
	}
}

func mention(username string) string {
	return "<@" + username + ">"
}

func (b *Bot) handleMessageEvent(msg slack.Msg) {
	b.log.WithFields(log.Fields{
		"Text":    msg.Text,
		"Channel": msg.Channel,
		"User":    msg.User,
	})

	tokens := strings.Split(msg.Text, " ")
	if len(tokens) == 0 {
		return
	}

	if tokens[0] == mention(username) {
		// Print help message
		if len(tokens) == 1 || tokens[1] == "help" {
			b.api.PostMessage(msg.Channel, helpMessage, noParams)
			return
		}

		if len(tokens) > 1 {
			member := model.Member{
				SlackID: msg.User,
			}

			if tokens[1] == "init" {
				if err := b.dal.CreateMember(&member); err != nil {
					b.log.WithError(err).Errorf("Error creating a new member with Slack ID %s", member.SlackID)
					b.api.PostMessage(msg.Channel, errorMessage, noParams)
					return
				}
				b.api.PostMessage(msg.Channel, "I've set up your profile! Please use these commands to add information:\n"+
					"`@rocket set name`\n`@rocket set email`\n`@rocket set github`\n`@rocket set program`", noParams)
				return
			}
			if err := b.dal.GetMemberBySlackID(&member); err != nil {
				b.log.WithError(err).Errorf("Error getting member by Slack ID %s", member.SlackID)
				b.api.PostMessage(msg.Channel, errorMessage, noParams)
				return
			}

			params := slack.PostMessageParameters{}

			if tokens[1] == "me" {
				params.Attachments = member.SlackAttachments()
				b.api.PostMessage(msg.Channel, "Your Launch Pad profile :rocket:", params)
				return
			}

			if len(tokens) > 3 {
				if tokens[1] == "set" {
					if tokens[2] == "name" {
						member.Name = strings.Join(tokens[3:], " ")
						if err := b.dal.SetMemberName(&member); err != nil {
							b.api.PostMessage(msg.Channel, errorMessage, noParams)
							return
						}
						params.Attachments = member.SlackAttachments()
						b.api.PostMessage(msg.Channel, "Your name has been updated! :simple_smile:", params)
						return
					}
					if tokens[2] == "email" {
						member.Email = tokens[3]
						if err := b.dal.SetMemberEmail(&member); err != nil {
							b.api.PostMessage(msg.Channel, errorMessage, noParams)
							return
						}
						params.Attachments = member.SlackAttachments()
						b.api.PostMessage(msg.Channel, "Your email has been updated! :simple_smile:", params)
						return
					}
					if tokens[2] == "github" {
						member.GithubUsername = tokens[3]
						if err := b.dal.SetMemberGitHubUsername(&member); err != nil {
							b.api.PostMessage(msg.Channel, errorMessage, noParams)
							return
						}
						params.Attachments = member.SlackAttachments()
						b.api.PostMessage(msg.Channel, "Your GitHub username has been updated! :simple_smile:", params)
						return
					}
					if tokens[2] == "program" {
						member.Program = tokens[3]
						if err := b.dal.SetMemberProgram(&member); err != nil {
							b.api.PostMessage(msg.Channel, errorMessage, noParams)
							return
						}
						params.Attachments = member.SlackAttachments()
						b.api.PostMessage(msg.Channel, "Your program has been updated! :simple_smile:", params)
						return
					}
				}
			}
		}
	}
}
