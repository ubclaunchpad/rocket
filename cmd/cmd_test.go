package cmd

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
	"github.com/ubclaunchpad/rocket/model"
)

var (
	anyString, _ = regexp.Compile(".*")
)

func getTestCommand(ch CommandHandler) *Command {
	return &Command{
		Name:     "test",
		HelpText: "fake command with two options",
		Options: map[string]*Option{
			"required": &Option{
				Key:      "required",
				HelpText: "this is a required option",
				Format:   anyString,
				Required: true,
			},
			"optional": &Option{
				Key:      "optional",
				HelpText: "this is an optional option",
				Format:   anyString,
				Required: false,
			},
		},
		HandleFunc: ch,
	}
}

func getTestContext(text string) Context {
	return Context{
		Message: &slack.Msg{
			Text: text,
		},
		User: model.Member{},
	}
}

func testHandler(context Context) (string, slack.PostMessageParameters) {
	return context.Message.Text, slack.PostMessageParameters{}
}

func TestCommand(t *testing.T) {
	ctx := getTestContext("@rocket test required=`gre at` optional=`awes =ome`")
	ch := func(c Context) (string, slack.PostMessageParameters) {
		ctx = c
		return "", slack.PostMessageParameters{}
	}
	cmd := getTestCommand(ch)
	_, _, err := cmd.Execute(ctx)
	assert.Equal(t, ctx.Options["required"].Value, "gre at")
	assert.Equal(t, ctx.Options["optional"].Value, "awes =ome")
	assert.Nil(t, err)
}

func TestInvalidCommand(t *testing.T) {
	ctx := getTestContext("@rocket ayyy required=`gre at`")
	cmd := getTestCommand(testHandler)
	_, _, err := cmd.Execute(ctx)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Invalid command"))
}

func TestCommandMissingRequiredOption(t *testing.T) {
	ctx := getTestContext("@rocket test optional=`noooo`")
	cmd := getTestCommand(testHandler)
	_, _, err := cmd.Execute(ctx)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Missing value for required option"))
}

func TestCommandDuplicateOption(t *testing.T) {
	ctx := getTestContext("@rocket test required=`ayy` required=`letsgo`")
	cmd := getTestCommand(testHandler)
	_, _, err := cmd.Execute(ctx)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Duplicate option"))
}

func TestCommandUnrecognizedOption(t *testing.T) {
	ctx := getTestContext("@rocket test plx=`plox`")
	cmd := getTestCommand(testHandler)
	_, _, err := cmd.Execute(ctx)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Unrecognized option"))
}

func TestCommandInvalidOptFormat(t *testing.T) {
	ctx := getTestContext("@rocket test required=`test`")
	cmd := getTestCommand(testHandler)
	cmd.Options["required"].Format, _ = regexp.Compile("^[0-9]")
	_, _, err := cmd.Execute(ctx)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Invalid format for option"))
}

func TestCommandHelp(t *testing.T) {
	cmd := getTestCommand(testHandler)
	res, _ := cmd.Help()
	assert.Equal(t, res,
		"Usage: @rocket test OPTIONS\n\nfake command with two options")
}
