package core

import (
	"fmt"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/bot"
	"github.com/ubclaunchpad/rocket/cmd"
)

// NewSetCmd returns a set command that sets user information
func NewSetCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:     "set",
		HelpText: "Set properties on your Launch Pad profile to a new values",
		Options: map[string]*cmd.Option{
			"name": &cmd.Option{
				Key:      "name",
				HelpText: "your full name",
				Format:   cmd.NameRegex,
				Required: false,
			},
			"email": &cmd.Option{
				Key:      "email",
				HelpText: "your email address",
				Format:   cmd.EmailRegex,
				Required: false,
			},
			"position": &cmd.Option{
				Key:      "position",
				HelpText: "your creative Launch Pad title",
				Format:   cmd.AnyRegex,
				Required: false,
			},
			"github": &cmd.Option{
				Key:      "github",
				HelpText: "your Github username",
				Format:   cmd.AnyRegex,
				Required: false,
			},
			"major": &cmd.Option{
				Key:      "major",
				HelpText: "your major at UBC",
				Format:   cmd.AnyRegex,
				Required: false,
			},
			"biography": &cmd.Option{
				Key:      "biography",
				HelpText: "a little bit about yourself (600 characters max)",
				Format:   cmd.AnyRegex,
				Required: false,
			},
		},
		HandleFunc: ch,
	}
}

// Generic command for setting some information about the sender's profile.
func (core *CorePlugin) set(c cmd.Context) (string, slack.PostMessageParameters) {
	params := slack.PostMessageParameters{}
	githubChanged := false

	if c.Options["name"].Value != "" {
		c.User.Name = c.Options["name"].Value
		if err := core.Bot.DAL.SetMemberName(&c.User); err != nil {
			log.WithError(err).Errorf("Failed to set name: %s", c.User.Name)
			return "Failed to set name " + c.User.Name, params
		}
	}

	if c.Options["email"].Value != "" {
		c.User.Email = c.Options["email"].Value
		if err := core.Bot.DAL.SetMemberEmail(&c.User); err != nil {
			log.WithError(err).Errorf("Failed to set email: %s", c.User.Email)
			return "Failed to set email " + c.User.Email, params
		}
	}

	if c.Options["github"].Value != "" {
		c.User.GithubUsername = c.Options["github"].Value
		// Check that the user exists
		exists, err := core.Bot.GitHub.UserExists(c.User.GithubUsername)
		if err != nil {
			log.WithError(err).Errorf("Error checking whether user %s exists", c.User.GithubUsername)
			return "Error checking whether user exists", params
		} else if !exists {
			return fmt.Sprintf("Github user %s does not exist", c.User.GithubUsername), params
		}

		// Add the user to our GitHub org by adding to `all` team
		if err := core.Bot.GitHub.AddUserToTeam(c.User.GithubUsername, bot.GithubAllTeamID); err != nil {
			log.WithError(err).Errorf("Failed to add %s to Launch Pad Github organization",
				c.User.GithubUsername)
			return "Failed to add you to Launch Pad's GitHub organization", params
		}

		// Finally, set their username in the DB
		if err := core.Bot.DAL.SetMemberGitHubUsername(&c.User); err != nil {
			log.WithError(err).Errorf("Failed to set GitHub username")
			return "Failed to set GitHub username", params
		}
		githubChanged = true
	}

	if c.Options["major"].Value != "" {
		c.User.Major = c.Options["major"].Value
		if err := core.Bot.DAL.SetMemberMajor(&c.User); err != nil {
			log.WithError(err).Error("Failed to set major")
			return "Failed to set major", params
		}
	}

	if c.Options["position"].Value != "" {
		c.User.Position = c.Options["position"].Value
		if err := core.Bot.DAL.SetMemberPosition(&c.User); err != nil {
			log.WithError(err).Error("Failed to set position")
			return "Failed to set position", params
		}
	}

	if c.Options["biography"].Value != "" {
		c.User.Biography = c.Options["biography"].Value
		// Max bio length is 600 characters
		if len(c.User.Biography) > 600 {
			return "Sorry, your biography must be at most 600 characters in length", params
		}
		if err := core.Bot.DAL.SetMemberBiography(&c.User); err != nil {
			log.WithError(err).Error("Failed to set biography")
			return "Failed to set biography", params
		}
	}

	params.Attachments = c.User.SlackAttachments()
	msg := "Your information has been updated :simple_smile:"
	if githubChanged {
		msg += "\nYou've also been added to our organization on GitHub, " +
			"so check your email for the invitation!"
	}
	return msg, params
}
