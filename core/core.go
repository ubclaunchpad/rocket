package core

import (
	"github.com/ubclaunchpad/rocket/bot"
	"github.com/ubclaunchpad/rocket/cmd"
)

// RocketPlugin stores the values required for accessing GitHub, Slack, Postgres,
// and Rocket's HTTP request handlers.
type RocketPlugin struct {
	Bot *bot.Bot
}

// New returns a new instance of the CorePlugin with the given bot.
func New(b *bot.Bot) *RocketPlugin {
	return &RocketPlugin{
		Bot: b,
	}
}

// Start initializes the pluin with the values it needs to do its job.
func (cp *RocketPlugin) Start() error {
	cp.Bot.Log.Info("Running CorePlugin")
	return nil
}

// Commands returns a list of commands this plugin makes available to the Bot.
func (cp *RocketPlugin) Commands() []*cmd.Command {
	return []*cmd.Command{
		NewHelpCmd(cp.help),
		NewSetCmd(cp.set),
		NewEditUserCmd(cp.editUser),
		NewViewUserCmd(cp.viewUser),
		NewViewTeamCmd(cp.viewTeam),
		NewAddUserCmd(cp.addUser),
		NewAddTeamCmd(cp.addTeam),
		NewEditTeamCmd(cp.editTeam),
		NewAddAdminCmd(cp.addAdmin),
		NewRemoveAdminCmd(cp.removeAdmin),
		NewRemoveUserCmd(cp.removeUser),
		NewRemoveTeamCmd(cp.removeTeam),
		NewTeamsCmd(cp.listTeams),
		NewAdminsCmd(cp.listAdmins),
		NewRefreshCmd(cp.refresh),
	}
}

// EventHandlers returns a mapping from Slack event name to event handler.
func (cp *RocketPlugin) EventHandlers() map[string]bot.EventHandler {
	// This plugin currently has no custom event handlers.
	return map[string]bot.EventHandler{}
}
