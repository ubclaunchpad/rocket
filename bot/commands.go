package bot

import (
	"github.com/ubclaunchpad/rocket/cmd"
)

var (
	// Commands is a list of all commands Rocket supports
	Commands = []*cmd.Command{HelpCmd, SetCmd, AddUserCmd,
		AddTeamCmd, AddAdminCmd, RemoveAdminCmd, RemoveTeamCmd, RemoveUserCmd,
		ViewTeamCmd, ViewUserCmd}

	// HelpCmd presents helpful information about Rocket commands
	HelpCmd = &cmd.Command{
		Name:     "help",
		HelpText: "Get help using Rocket commands",
		Options: map[string]*cmd.Option{
			"command": &cmd.Option{
				Key:      "command",
				HelpText: "get help using a particular Rocket command",
				Format:   alphaRegex,
			},
		},
		Args: []cmd.Argument{},
	}
	// SetCmd sets user information
	SetCmd = &cmd.Command{
		Name:     "set",
		HelpText: "Set properties on your Launch Pad profile to a new values",
		Options: map[string]*cmd.Option{
			"name": &cmd.Option{
				Key:      "name",
				HelpText: "your full name",
				Format:   nameRegex,
			},
			"email": &cmd.Option{
				Key:      "email",
				HelpText: "your email address",
				Format:   emailRegex,
			},
			"position": &cmd.Option{
				Key:      "position",
				HelpText: "your creative Launch Pad title",
				Format:   anyRegex,
			},
			"github": &cmd.Option{
				Key:      "github",
				HelpText: "your Github username",
				Format:   anyRegex,
			},
			"major": &cmd.Option{
				Key:      "major",
				HelpText: "your major at UBC",
				Format:   anyRegex,
			},
		},
		Args: []cmd.Argument{},
	}
	// AddUserCmd adds a user
	AddUserCmd = &cmd.Command{
		Name:     "add-user",
		HelpText: "Add a user to a team",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:      "member",
				HelpText:  "the Slack handle of the user to add to a team",
				Format:    anyRegex,
				MultiWord: false,
			},
			cmd.Argument{
				Name:      "team-name",
				HelpText:  "the team to add the user to",
				Format:    anyRegex,
				MultiWord: true,
			},
		},
	}
	// AddAdminCmd makes an existing user an admin (this action can only be
	// performed by admins)
	AddAdminCmd = &cmd.Command{
		Name:     "add-admin",
		HelpText: "Make an existing user an admin (admins only)",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:      "username",
				HelpText:  "the Slack handle of the user to make an admin",
				Format:    anyRegex,
				MultiWord: false,
			},
		},
	}
	// AddTeamCmd creats a new Launch Pad team
	AddTeamCmd = &cmd.Command{
		Name:     "add-team",
		HelpText: "Create a new Launch Pad team",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:      "team-name",
				HelpText:  "the name of the new team",
				Format:    anyRegex,
				MultiWord: true,
			},
		},
	}
	// RemoveUserCmd removes a user
	RemoveUserCmd = &cmd.Command{
		Name:     "remove-user",
		HelpText: "Remove a user from a team",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:     "username",
				HelpText: "the Slack handle of the user to remove from a team",
				Format:   anyRegex,
			},
			cmd.Argument{
				Name:      "team",
				HelpText:  "the team to remove the user from",
				Format:    anyRegex,
				MultiWord: true,
			},
		},
	}
	// RemoveAdminCmd makes an existing user an admin (this action can only be
	// performed by admins)
	RemoveAdminCmd = &cmd.Command{
		Name:     "remove-admin",
		HelpText: "Remove admin rights from a user (admins only)",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:      "username",
				HelpText:  "the Slack handle of the user to remove admin rights from",
				Format:    anyRegex,
				MultiWord: false,
			},
		},
	}
	// RemoveTeamCmd removes a new Launch Pad team
	RemoveTeamCmd = &cmd.Command{
		Name:     "remove-team",
		HelpText: "Delete a new Launch Pad team",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:      "team-name",
				HelpText:  "the name of the team to remove",
				Format:    anyRegex,
				MultiWord: true,
			},
		},
	}
	// ViewUserCmd displays information about a user
	ViewUserCmd = &cmd.Command{
		Name:     "view-user",
		HelpText: "View information about a user",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:      "username",
				HelpText:  "the slack handle of the user to view",
				Format:    anyRegex,
				MultiWord: false,
			},
		},
	}
	// ViewTeamCmd displays information about a user
	ViewTeamCmd = &cmd.Command{
		Name:     "view-team",
		HelpText: "View information about a Launch Pad team",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:      "team-name",
				HelpText:  "the name of the team to view",
				Format:    anyRegex,
				MultiWord: true,
			},
		},
	}
)
