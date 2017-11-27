package bot

import (
	"strings"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/config"
	"github.com/ubclaunchpad/rocket/data"
	"github.com/ubclaunchpad/rocket/github"
	"github.com/ubclaunchpad/rocket/model"
)

const (
	// Our Slack Bot's username on the UBC Launch Pad Slack
	username = "U5RU9TB38"

	commandsMessage = "```\n@rocket set name <name>\n@rocket set email <email>\n" +
		"@rocket set github <username>\n@rocket set major <major>\n@rocket set position <position>```"

	helpMessage = "Hi there, I'm Rocket, Launch Pad's friendly neighbourhood Slack bot! :rocket:\n" +
		"You view your profile with `@rocket me`.\n" +
		"You can update your profile too!\n" +
		commandsMessage

	// Default message to send when any error occurs
	errorMessage = "Oops, an error occurred :robot_face:. Bruno must have coded a bug... Sorry about that!"

	// ID for the `all` team that everyone should be on
	githubAllTeamID = 2467607
)

var noParams = slack.PostMessageParameters{}

// CommandContext contains the Slack message, command arguments, and sender
// for commands received through Slack.
type CommandContext struct {
	msg  *slack.Msg
	args []string
	user model.Member
}

// CommandHandler defines an interface for functions that respond to commands
// should take.
type CommandHandler func(*CommandContext)

// Bot represents an instance of the Rocket Slack bot. Only one should be
// created under normal circumstances.
type Bot struct {
	token    string
	api      *slack.Client
	rtm      *slack.RTM
	dal      *data.DAL
	gh       *github.API
	log      *log.Entry
	commands map[string]CommandHandler
	users    map[string]slack.User
}

// New constructs and returns a new Slack bot instance. It creates a new RTM
// object to receive incoming messages, populates a cache with users, and
// sets up command handlers.
func New(cfg *config.Config, dal *data.DAL, gh *github.API, log *log.Entry) *Bot {
	api := slack.New(cfg.SlackToken)

	b := &Bot{
		token: cfg.SlackToken,
		api:   api,
		rtm:   api.NewRTM(),
		dal:   dal,
		gh:    gh,
		log:   log,
	}

	commands := map[string]CommandHandler{
		"help":    b.help,
		"me":      b.me,
		"set":     b.set,
		"add":     b.add,
		"remove":  b.remove,
		"view":    b.view,
		"refresh": b.refresh,
	}
	b.commands = commands

	b.PopulateUsers()

	return b
}

// Start causes an already initialized bot instance to begin listening for
// and responding to commands sent on its Slack channel.
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

// PopulateUsers retrieves list of users from API and populates the bot
// instance's cache.
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

// SendErrorMessage sends a generic error message back to the sender and
// logs the specific error that occurred.
func (b *Bot) SendErrorMessage(channel string, err error, msg string) {
	errorMsg := errorMessage
	if len(msg) > 0 {
		errorMsg = msg
	}
	b.api.PostMessage(channel, errorMsg, noParams)
	b.log.WithError(err).Error(msg)
}

// Generic handler for any new message we receive. Determines whether the
// message is meant to be a command (if we need to take action for it),
// populates the command context object for the message, and calls the
// appropriate handler.
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

	// Create member if doesn't already exist (this acts like an upsert)
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

	// Retrieves the full member object from the database
	if err := b.dal.GetMemberBySlackID(&member); err != nil {
		b.log.WithError(err).Errorf("Error getting member by Slack ID %s", member.SlackID)
		b.api.PostMessage(msg.Channel, errorMessage, noParams)
		return
	}

	args := strings.Fields(msg.Text)
	if len(args) == 0 {
		return
	}

	// A command is defined by being prefixed by our username
	// i.e. "@rocket <command> <arg1> ..."
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

// Handler for when a user changes their profile, or a user is added/deleted.
// Creates the member if they don't already exist and sets their profile image.
func (b *Bot) handleUserChange(user slack.User) {
	b.users[user.ID] = user

	member := model.Member{
		SlackID:  user.ID,
		ImageURL: user.Profile.Image192,
	}

	// Create user if doesn't exist
	if err := b.dal.CreateMember(&member); err != nil {
		b.log.WithError(err).Errorf("Error creating user with Slack ID %s", member.SlackID)
	}

	// Update image URL
	if err := b.dal.SetMemberImageURL(&member); err != nil {
		b.log.WithError(err).Errorf("Error setting image URL for Slack ID %s", member.SlackID)
	}
}
