package bot

import (
	"strings"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
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

func getTestBot() *Bot {
	b := &Bot{}
	b.commands = map[string]*cmd.Command{
		"help":        NewHelpCmd(b.help),
		"set":         NewSetCmd(b.set),
		"viewuser":    NewViewUserCmd(b.viewUser),
		"viewteam":    NewViewTeamCmd(b.viewTeam),
		"adduser":     NewAddUserCmd(b.addUser),
		"addteam":     NewAddTeamCmd(b.addTeam),
		"addadmin":    NewAddAdminCmd(b.addAdmin),
		"removeadmin": NewRemoveAdminCmd(b.removeAdmin),
		"removeuser":  NewRemoveUserCmd(b.removeUser),
		"removeteam":  NewRemoveTeamCmd(b.removeTeam),
		"teams":       NewTeamsCmd(b.listTeams),
		"refresh":     NewRefreshCmd(b.refresh),
	}
	return b
}

func TestHelp(t *testing.T) {
	ctx := getTestContext("@rocket help")
	b := getTestBot()
	res, _, err := b.commands["help"].Execute(ctx)
	t.Log(res)
	assert.Nil(t, err)
}

func TestHelpWithCommand(t *testing.T) {
	b := getTestBot()
	for _, cmd := range b.commands {
		text := "@rocket help command=`" + cmd.Name + "`"
		ctx := getTestContext(text)
		res, _, err := b.commands["help"].Execute(ctx)
		t.Log(res)
		assert.Nil(t, err)
		assert.True(t, strings.Contains(res, "Usage:"))
	}
}

func TestHelpWithInvalidCommand(t *testing.T) {
	text := "@rocket help command=`blabla`"
	ctx := getTestContext(text)
	b := getTestBot()
	res, _, err := b.commands["help"].Execute(ctx)
	t.Log(res)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(res, "is not a Rocket command"))
}
