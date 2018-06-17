package core

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

// NewTechLeadsCmd returns a teams command that displays a list of Launch Pad teams
func NewTechLeadsCmd(ch cmd.CommandHandler) *cmd.Command {
	return &cmd.Command{
		Name:       "tech-leads",
		HelpText:   "List Launch Pad tech leads",
		Options:    map[string]*cmd.Option{},
		HandleFunc: ch,
	}
}

// listAdmins displays Launch Pad admins
func (core *Plugin) listTechLeads(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	members := model.Members{}
	if err := core.Bot.DAL.GetTechLeads(&members); err != nil {
		log.WithError(err).Error("failed to get tech leads")
		return "Failed to get tech leads", noParams
	}
	names := ""
	for _, member := range members {
		names += member.Name + "\n"
	}
	return names, noParams
}
