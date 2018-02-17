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

func TestHelp(t *testing.T) {
	ctx := getTestContext("@rocket help")
	b := Bot{}
	HelpCmd.HandleFunc = b.help
	res, _, err := HelpCmd.Execute(ctx)
	t.Log(res)
	assert.Nil(t, err)
}

func TestHelpWithCommand(t *testing.T) {
	for _, cmd := range Commands {
		text := "@rocket help --command=`" + cmd.Name + "`"
		ctx := getTestContext(text)
		b := Bot{}
		HelpCmd.HandleFunc = b.help
		res, _, err := HelpCmd.Execute(ctx)
		t.Log(res)
		assert.Nil(t, err)
		assert.True(t, strings.Contains(res, "Usage:"))
	}
}

func TestHelpWithInvalidCommand(t *testing.T) {
	text := "@rocket help --command=`blabla`"
	ctx := getTestContext(text)
	b := Bot{}
	HelpCmd.HandleFunc = b.help
	res, _, err := HelpCmd.Execute(ctx)
	t.Log(res)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(res, "is not a Rocket command"))
}
