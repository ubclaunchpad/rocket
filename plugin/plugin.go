package plugin

import (
	"github.com/ubclaunchpad/rocket/bot"
	"github.com/ubclaunchpad/rocket/cmd"

	// Place registered plugin imports here
	"github.com/ubclaunchpad/rocket/plugins/core"
	"github.com/ubclaunchpad/rocket/plugins/welcome"
)

// Plugin is any type that exposes Slack commands and event handlers, and can
// be started.
type Plugin interface {
	// Starts the plugin or returns an error if one occurred.
	// Use this as an opportnity to start background goroutines or do any other
	// additional setup for your plugin.
	Start() error
	// Returns a slice of commands that the plugin handles.
	Commands() []*cmd.Command
	// Returns a mapping from event type to a event handler.
	// See https://api.slack.com/rtm for event types.
	EventHandlers() map[string]bot.EventHandler
}

// RegisterPlugins registers commands and event handlers from Rocket plugins
// and starts the plugins. Returns an error if a plugin could not be registered.
func RegisterPlugins(b *bot.Bot) error {
	// Add your plugin to this list
	plugins := []Plugin{
		core.New(b),
		welcome.New(b),
	}
	for _, p := range plugins {
		if err := registerPlugin(p, b); err != nil {
			return err
		}
	}
	return nil
}

// RegisterPlugins registers commands and event handlers from the given plugin
// and starts the plugin. Returns an error if a plugin could not be registered.
func registerPlugin(p Plugin, b *bot.Bot) error {
	if err := b.RegisterCommands(p.Commands()); err != nil {
		return err
	}
	b.RegisterEventHandlers(p.EventHandlers())
	p.Start()
	return nil
}
