package plugin

import (
	"github.com/ubclaunchpad/rocket/bot"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/core"
)

// Plugin is any type that exposes Slack commands and event handlers, and can
// be started.
type Plugin interface {
	// Starts the plugin.
	// Use this as an opportnity to start background goroutines or do any other
	// additional setup for your plugin.
	Start()
	// Returns a slice of commands that the plugin handles.
	Commands() []*cmd.Command
	// Returns a mapping from event type to a event handler.
	// See https://api.slack.com/rtm for event types.
	EventHandlers() map[string]bot.EventHandler
}

// RegisterPlugins registers commands and event handlers from Rocket plugins
// and starts the plugins.
func RegisterPlugins(b *bot.Bot) {
	// Add your new plugins here
	registerPlugin(core.New(b), b)
}

// RegisterPlugins registers commands and event handlers from the given plugin
// and starts the plugin.
func registerPlugin(p Plugin, b *bot.Bot) {
	b.RegisterCommands(p.Commands())
	b.RegisterEventHandlers(p.EventHandlers())
	p.Start()
}
