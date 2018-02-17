package bot

import (
	"github.com/ubclaunchpad/rocket/cmd"
)

var (
	// Commands is a list of all commands Rocket supports
	Commands = []*cmd.Command{HelpCmd, SetCmd, AddUserCmd,
		AddTeamCmd, AddAdminCmd}

	// HelpCmd presents helpful information about Rocket commands
	HelpCmd = &cmd.Command{
		Name:     "help",
		HelpText: "get help using Rocket commands",
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
		HelpText: "set properties on your Launch Pad profile to a new values",
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
		HelpText: "add a user to the Launch Pad organization",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:     "member",
				HelpText: "the Slack handle of the user to add to a team",
				Format:   usernameRegex,
			},
			cmd.Argument{
				Name:     "team-name",
				HelpText: "the team to add the user to",
				Format:   anyRegex,
			},
		},
	}
	// AddAdminCmd makes an existing user an admin (this action can only be
	// performed by admins)
	AddAdminCmd = &cmd.Command{
		Name:     "add-admin",
		HelpText: "make an existing user an admin (can only be performed by admins)",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:     "username",
				HelpText: "the Slack handle of the user to make an admin",
				Format:   usernameRegex,
			},
		},
	}
	// AddTeamCmd creats a new Launch Pad team
	AddTeamCmd = &cmd.Command{
		Name:     "add-team",
		HelpText: "create a new Launch Pad team",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:     "team-name",
				HelpText: "the name of the new team",
				Format:   anyRegex,
			},
		},
	}
	// RemoveUserCmd removes a user
	RemoveUserCmd = &cmd.Command{
		Name:     "remove-user",
		HelpText: "remove a user to the Launch Pad organization",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:     "member",
				HelpText: "the Slack handle of the user to remove from a team",
				Format:   usernameRegex,
			},
			cmd.Argument{
				Name:     "team",
				HelpText: "the team to remove the user from",
				Format:   anyRegex,
			},
		},
	}
	// RemoveAdminCmd makes an existing user an admin (this action can only be
	// performed by admins)
	RemoveAdminCmd = &cmd.Command{
		Name:     "remove-admin",
		HelpText: "remove admin rights from a user (can only be performed by admins)",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:     "username",
				HelpText: "the Slack handle of the user to remove admin rights from",
				Format:   usernameRegex,
			},
		},
	}
	// RemoveTeamCmd removes a new Launch Pad team
	RemoveTeamCmd = &cmd.Command{
		Name:     "remove-team",
		HelpText: "delete a new Launch Pad team",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:     "team-name",
				HelpText: "the name of the team to remove",
				Format:   anyRegex,
			},
		},
	}
	// ViewUserCmd displays information about a user
	ViewUserCmd = &cmd.Command{
		Name:     "view-user",
		HelpText: "view information about a user",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:     "username",
				HelpText: "the slack handle of the user to view",
				Format:   usernameRegex,
			},
		},
	}
	// ViewTeamCmd displays information about a user
	ViewTeamCmd = &cmd.Command{
		Name:     "view-team",
		HelpText: "view information about a Launch Pad team",
		Options:  map[string]*cmd.Option{},
		Args: []cmd.Argument{
			cmd.Argument{
				Name:     "team-name",
				HelpText: "the name of the team to view",
				Format:   anyRegex,
			},
		},
	}
)
