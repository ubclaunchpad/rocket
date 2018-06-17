package core

import (
	"strings"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
	"github.com/ubclaunchpad/rocket/bot"
	"github.com/ubclaunchpad/rocket/cmd"
	"github.com/ubclaunchpad/rocket/model"
)

func getTestContext(text string) cmd.Context {
	return cmd.Context{
		Message: &slack.Msg{
			Text: text,
		},
		User: model.Member{
			Name:           "Rocket Man",
			Email:          "rocketman@test.mars",
			GithubUsername: "rocketman",
			Major:          "Rocket Stuff",
			Position:       "On a rocket",
		},
	}
}

func getTestBot() *bot.Bot {
	b := &bot.Bot{}
	cp := New(b)
	b.Commands = map[string]*cmd.Command{
		"help":        NewHelpCmd(cp.help),
		"set":         NewSetCmd(cp.set),
		"viewuser":    NewViewUserCmd(cp.viewUser),
		"viewteam":    NewViewTeamCmd(cp.viewTeam),
		"adduser":     NewAddUserCmd(cp.addUser),
		"addteam":     NewAddTeamCmd(cp.addTeam),
		"setadmin":    NewToggleAdminCmd(cp.toggleAdmin),
		"removeuser":  NewRemoveUserCmd(cp.removeUser),
		"removeteam":  NewRemoveTeamCmd(cp.removeTeam),
		"teams":       NewTeamsCmd(cp.listTeams),
		"refresh":     NewRefreshCmd(cp.refresh),
		"settechlead": NewToggleTechLeadCmd(cp.toggleTechLead),
		"techleads":   NewTechLeadsCmd(cp.listTechLeads),
	}
	return b
}

func TestHelp(t *testing.T) {
	ctx := getTestContext("@rocket help")
	b := getTestBot()
	res, _, err := b.Commands["help"].Execute(ctx)
	t.Log(res)
	assert.Nil(t, err)
}

func TestHelpWithCommand(t *testing.T) {
	b := getTestBot()
	for _, cmd := range b.Commands {
		text := "@rocket help command={" + cmd.Name + "}"
		ctx := getTestContext(text)
		res, _, err := b.Commands["help"].Execute(ctx)
		t.Log(res)
		assert.Nil(t, err)
		assert.True(t, strings.Contains(res, "Usage:"))
	}
}

func TestHelpWithInvalidCommand(t *testing.T) {
	text := "@rocket help command={blabla}"
	ctx := getTestContext(text)
	b := getTestBot()
	res, _, err := b.Commands["help"].Execute(ctx)
	t.Log(res)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(res, "is not a Rocket command"))
}
