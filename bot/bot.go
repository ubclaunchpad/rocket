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

type CommandContext struct {
	msg  *slack.Msg
	args []string
	user model.Member
}

type CommandHandler func(*CommandContext)

type Bot struct {
	token    string
	api      *slack.Client
	rtm      *slack.RTM
	dal      *data.DAL
	log      *log.Entry
	commands map[string]CommandHandler
	users    map[string]slack.User
}

func New(cfg *config.Config, dal *data.DAL, log *log.Entry) *Bot {
	api := slack.New(cfg.SlackToken)

	b := &Bot{
		token: cfg.SlackToken,
		api:   api,
		rtm:   api.NewRTM(),
		dal:   dal,
		log:   log,
	}

	commands := map[string]CommandHandler{
		"help":    b.help,
		"me":      b.me,
		"set":     b.set,
		"add":     b.add,
		"remove":  b.remove,
		"refresh": b.refresh,
	}
	b.commands = commands

	b.PopulateUsers()

	return b
}

func (b *Bot) Start() {
	go b.rtm.ManageConnection()

	for evt := range b.rtm.IncomingEvents {
		switch evt.Data.(type) {
		// Check for and respond to commands
		case *slack.MessageEvent:
			b.handleMessageEvent(evt.Data.(*slack.MessageEvent).Msg)
		// Update our user cache when new member joins or user data changes
		case *slack.TeamJoinEvent:
			b.handleUserChange(evt.Data.(*slack.TeamJoinEvent).User)
		case *slack.UserChangeEvent:
			b.handleUserChange(evt.Data.(*slack.UserChangeEvent).User)
		}
	}
}

// PopulateUsers retrieves list of users from API and populates in-memory cache
func (b *Bot) PopulateUsers() {
	users, err := b.api.GetUsers()
	if err != nil {
		b.log.WithError(err).Error("Failed to populate users")
	}
	b.users = make(map[string]slack.User)
	for _, u := range users {
		b.users[u.ID] = u
	}
}

func (b *Bot) SendErrorMessage(channel string, err error, msg string) {
	errorMsg := errorMessage
	if len(msg) > 0 {
		errorMsg = msg
	}
	b.api.PostMessage(channel, errorMsg, noParams)
	b.log.WithError(err).Error(msg)
}

func (b *Bot) handleMessageEvent(msg slack.Msg) {
	b.log.WithFields(log.Fields{
		"Text":    msg.Text,
		"Channel": msg.Channel,
		"User":    msg.User,
	})

	// Ignore messages from bots
	if len(msg.User) == 0 {
		return
	}

	member := model.Member{
		SlackID:  msg.User,
		ImageURL: b.users[msg.User].Profile.Image192,
	}

	// Create member if doesn't already exist
	if err := b.dal.CreateMember(&member); err != nil {
		b.log.WithError(err).Errorf("Error creating member with Slack ID %s", member.SlackID)
		b.api.PostMessage(msg.Channel, errorMessage, noParams)
		return
	}

	// Set member image to their slack profile image
	if err := b.dal.SetMemberImageURL(&member); err != nil {
		b.log.WithError(err).Errorf("Error setting member image URL")
		b.api.PostMessage(msg.Channel, errorMessage, noParams)
		return
	}

	if err := b.dal.GetMemberBySlackID(&member); err != nil {
		b.log.WithError(err).Errorf("Error getting member by Slack ID %s", member.SlackID)
		b.api.PostMessage(msg.Channel, errorMessage, noParams)
		return
	}

	args := strings.Fields(msg.Text)
	if len(args) == 0 {
		return
	}

	// Command message
	if args[0] == toMention(username) {
		context := &CommandContext{
			msg:  &msg,
			args: args[1:],
			user: member,
		}

		if len(args) > 1 {
			command := args[1]
			handler, ok := b.commands[command]
			if !ok {
				handler = b.help
			}
			handler(context)
		} else {
			b.help(context)
		}
	}
}

func (b *Bot) handleUserChange(user slack.User) {
	b.users[user.ID] = user
}
