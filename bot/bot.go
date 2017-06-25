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

	commandsMessage = "```\n@rocket set name <name>\n@rocket set email <email>\n" +
		"@rocket set github <username>\n@rocket set major <major>\n@rocket set position <position>```"

	helpMessage = "Hi there, I'm Rocket, Launch Pad's friendly neighbourhood Slack bot! :rocket:\n" +
		"You view your profile with `@rocket me`.\n" +
		"You can update your profile too!\n" +
		commandsMessage

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

	if tokens[0] == toMention(username) {
		// Print help message
		if len(tokens) == 1 || tokens[1] == "help" {
			b.api.PostMessage(msg.Channel, helpMessage, noParams)
			return
		}

		if len(tokens) > 1 {
			member := model.Member{
				SlackID: msg.User,
			}

			// Create member if doesn't already exist
			if err := b.dal.CreateMember(&member); err != nil {
				b.log.WithError(err).Errorf("Error creating member with Slack ID %s", member.SlackID)
				b.api.PostMessage(msg.Channel, errorMessage, noParams)
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
						member.Email = parseEmail(tokens[3])
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
					if tokens[2] == "major" {
						member.Major = tokens[3]
						if err := b.dal.SetMemberMajor(&member); err != nil {
							b.api.PostMessage(msg.Channel, errorMessage, noParams)
							return
						}
						params.Attachments = member.SlackAttachments()
						b.api.PostMessage(msg.Channel, "Your major has been updated! :simple_smile:", params)
						return
					}
					if tokens[2] == "position" {
						member.Position = strings.Join(tokens[3:], " ")
						if err := b.dal.SetMemberPosition(&member); err != nil {
							b.api.PostMessage(msg.Channel, errorMessage, noParams)
							return
						}
						params.Attachments = member.SlackAttachments()
						b.api.PostMessage(msg.Channel, "You position has been updated! :simple_smile:", params)
					}
				}
			}
		}
	}
}
