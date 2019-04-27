package core

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewAdminsCmd returns a teams command that displays a list of Launch Pad teams
func NewAdminsCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:       "admins",
		HelpText:   "List Launch Pad admins",
		Options:    map[string]*cmd.Option{},
		HandleFunc: ch,
	}
}

// listAdmins displays Launch Pad admins
func (core *Plugin) listAdmins(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	members := model.Members{}
	if err := core.Bot.DAL.GetAdmins(&members); err != nil {
		log.WithError(err).Error("failed to get admins")
		return "Failed to get admins", noParams
	}
	msg := ""
	for _, member := range members {
		msg += member.Name + "\n"
	}
	if len(msg) == 0 {
		msg = "There are currently no admins :feelsbadman:"
	}
	return msg, noParams
}
