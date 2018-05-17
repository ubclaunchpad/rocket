package bot

import (
	"strings"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/config"
	"github.com/ubclaunchpad/rocket/data"
	"github.com/ubclaunchpad/rocket/github"
	"github.com/ubclaunchpad/rocket/model"
)

const (
	// Our Slack Bot's username on the UBC Launch Pad Slack
	username = "U5RU9TB38"

	// Default message to send when any error occurs
	errorMessage = "Oops, an error occurred :robot_face:. Bruno must have " +
		"coded a bug... Sorry about that!"

	// ID for the `all` team that everyone should be on
	GithubAllTeamID = 2467607
)

var noParams = slack.PostMessageParameters{}

// EventHandler is any function that handles a Slack event
type EventHandler func(slack.RTMEvent)

// Bot represents an instance of the Rocket Slack bot. Only one should be
// created under normal circumstances.
type Bot struct {
	token    string
	API      *slack.Client
	rtm      *slack.RTM
	DAL      *data.DAL
	GitHub   *github.API
	Log      *log.Entry
	Commands map[string]*cmd.Command
	handlers map[string][]EventHandler
	Users    map[string]slack.User
}

// New constructs and returns a new Slack bot instance. It creates a new RTM
// object to receive incoming messages, populates a cache with users, and
// sets up command handlers.
func New(cfg *config.Config, dal *data.DAL, gh *github.API, log *log.Entry) *Bot {
	api := slack.New(cfg.SlackToken)

	b := &Bot{
		token:    cfg.SlackToken,
		API:      api,
		rtm:      api.NewRTM(),
		DAL:      dal,
		GitHub:   gh,
		Log:      log,
		Commands: map[string]*cmd.Command{},
		handlers: map[string][]EventHandler{},
	}
	b.PopulateUsers()

	// Register default Slack event handlers
	b.RegisterEventHandlers(map[string]EventHandler{
		"message":     b.handleMessageEvent,
		"team_join":   b.handleUserChange,
		"user_change": b.handleUserChange,
	})
	return b
}

// RegisterEventHandlers registers a handlers for different events. These
// handlers will be called when an event of the corresponding type is received.
func (b *Bot) RegisterEventHandlers(handlers map[string]EventHandler) {
	for evt, handler := range handlers {
		if b.handlers[evt] == nil {
			b.handlers[evt] = []EventHandler{handler}
		} else {
			b.handlers[evt] = append(b.handlers[evt], handler)
		}
		b.Log.Infof("registered handler for %s Slack event", evt)
	}
}

// RegisterCommands registers commands that the bot should handle.
func (b *Bot) RegisterCommands(commands []*cmd.Command) {
	for _, c := range commands {
		if b.Commands[c.Name] != nil {
			b.Log.Errorf("attempt to register duplicate commands: %s", c.Name)
			continue
		}
		b.Commands[c.Name] = c
		b.Log.Infof("registered command %s", c.Name)
	}
}

// Start causes an already initialized bot instance to begin listening for
// and responding to commands sent on its Slack channel.
func (b *Bot) Start() {
	go b.rtm.ManageConnection()

	for evt := range b.rtm.IncomingEvents {
		// Call any registered event handlers that are expecting events of this
		// type.
		if handlers := b.handlers[evt.Type]; handlers != nil {
			for _, handler := range handlers {
				handler(evt)
			}
		}
	}
}

// PopulateUsers retrieves list of users from API, populates the bot
// instance's cache, and updates any member entries in the DB with any relevant
// info from their Slack profiles.
func (b *Bot) PopulateUsers() {
	users, err := b.API.GetUsers()
	if err != nil {
		b.Log.WithError(err).Error("Failed to populate users")
	}

	b.Users = make(map[string]slack.User)
	for _, u := range users {
		b.Users[u.ID] = u
		member := &model.Member{
			SlackID:  u.ID,
			Name:     u.Profile.RealName,
			IsAdmin:  u.IsAdmin,
			Email:    u.Profile.Email,
			Position: u.Profile.Title,
		}
		if err := b.DAL.UpdateMember(member); err != nil {
			b.Log.WithError(err).Error("failed to update member " + member.SlackID)
		}
		b.Log.Debug("Successfully updated user ", member.SlackID)
	}
}

// SendErrorMessage sends a generic error message back to the sender and
// logs the specific error that occurred.
func (b *Bot) SendErrorMessage(channel string, err error, msg string) {
	errorMsg := errorMessage
	if len(msg) > 0 {
		errorMsg = msg
	}
	b.API.PostMessage(channel, errorMsg, noParams)
	b.Log.WithError(err).Error(msg)
}

// handleMessageEvent is a generic handler for any new message we receive.
// Determines whether the message is meant to be a command (if we need to
// take action for it), populates the command context object for the message,
// and calls the appropriate handler.
func (b *Bot) handleMessageEvent(evt slack.RTMEvent) {
	msg := evt.Data.(*slack.MessageEvent).Msg
	b.Log.WithFields(log.Fields{
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
		ImageURL: b.Users[msg.User].Profile.Image192,
	}

	// Create member if doesn't already exist (this acts like an upsert)
	if err := b.DAL.CreateMember(&member); err != nil {
		b.Log.WithError(err).Errorf("Error creating member with Slack ID %s", member.SlackID)
		b.API.PostMessage(msg.Channel, errorMessage, noParams)
		return
	}

	// Set member image to their slack profile image
	if err := b.DAL.SetMemberImageURL(&member); err != nil {
		b.Log.WithError(err).Errorf("Error setting member image URL")
		b.API.PostMessage(msg.Channel, errorMessage, noParams)
		return
	}

	// Retrieves the full member object from the database
	if err := b.DAL.GetMemberBySlackID(&member); err != nil {
		b.Log.WithError(err).Errorf("Error getting member by Slack ID %s", member.SlackID)
		b.API.PostMessage(msg.Channel, errorMessage, noParams)
		return
	}

	args := strings.Fields(msg.Text)
	if len(args) == 0 {
		return
	}

	// A command is defined by being prefixed by our username
	// i.e. "@rocket <command> <arg1> ..."
	if args[0] == cmd.ToMention(username) {
		context := cmd.Context{
			Message: &msg,
			User:    member,
		}

		var cmd *cmd.Command
		if len(args) > 1 {
			command := args[1]
			cmd = b.Commands[command]
			if cmd == nil {
				cmd = b.Commands["help"]
			}
		} else {
			cmd = b.Commands["help"]
		}
		res, params, err := cmd.Execute(context)
		if err != nil {
			log.WithError(err).Error("Failed to execute command")
			b.SendErrorMessage(context.Message.Channel, err, err.Error())
		}
		b.API.PostMessage(context.Message.Channel, res, params)
	}
}

// Handler for when a user changes their profile, or a user is added/deleted.
// Creates the member if they don't already exist and sets their profile image.
func (b *Bot) handleUserChange(evt slack.RTMEvent) {
	var user slack.User

	// This function is only called for team join or user change events, so
	// check which case we are in before proceeding.
	if evt.Type == "team_join" {
		user = evt.Data.(*slack.TeamJoinEvent).User
	} else {
		user = evt.Data.(*slack.UserChangeEvent).User
	}

	b.Users[user.ID] = user
	member := model.Member{
		SlackID:  user.ID,
		ImageURL: user.Profile.Image192,
	}

	// Create user if doesn't exist
	if err := b.DAL.CreateMember(&member); err != nil {
		b.Log.WithError(err).Errorf("Error creating user with Slack ID %s", member.SlackID)
	}

	// Update image URL
	if err := b.DAL.SetMemberImageURL(&member); err != nil {
		b.Log.WithError(err).Errorf("Error setting image URL for Slack ID %s", member.SlackID)
	}
}
