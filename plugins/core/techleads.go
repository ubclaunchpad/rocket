package core

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/ubclaunchpad/rocket/cmd"
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

// listTechLeads displays Launch Pad tech leads
func (core *Plugin) listTechLeads(c cmd.Context) (string, slack.PostMessageParameters) {
	noParams := slack.PostMessageParameters{}
	members, err := core.Bot.DAL.GetTechLeads()
	if err != nil {
		log.WithError(err).Error("failed to get tech leads")
		return "Failed to get tech leads", noParams
	}
	msg := ""
	for _, member := range *members {
		msg += member.Name + "\n"
	}
	if len(msg) == 0 {
		msg = "There are currently no tech leads :feelsbadman:"
	}
	return msg, noParams
}
