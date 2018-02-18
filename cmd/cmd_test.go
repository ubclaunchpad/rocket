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
		Name:     "testcommand",
		HelpText: "lets go dude!",
		Options: map[string]*Option{
			"myopt": &Option{
				Key:      "myopt",
				HelpText: "what a sick option",
				Format:   anyString,
			},
		},
		Args: []Argument{
			Argument{
				Name:      "testarg",
				HelpText:  "what a cool arg",
				Format:    anyString,
				MultiWord: false,
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
	ctx := getTestContext("@rocket testcommand --myopt=`great` 2")
	ch := func(c Context) (string, slack.PostMessageParameters) {
		ctx = c
		return "", slack.PostMessageParameters{}
	}
	cmd := getTestCommand(ch)
	_, _, err := cmd.Execute(ctx)
	assert.Equal(t, ctx.Args[0].Value, "2")
	assert.Equal(t, ctx.Options["myopt"].Value, "great")
	assert.Nil(t, err)
}

func TestInvalidCommand(t *testing.T) {
	ctx := getTestContext("@rocket ayyy --myopt=`great` 2")
	cmd := getTestCommand(testHandler)
	_, _, err := cmd.Execute(ctx)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Invalid command"))
}

func TestCommandMissingArg(t *testing.T) {
	ctx := getTestContext("@rocket testcommand --myopt=`true`")
	cmd := getTestCommand(testHandler)
	_, _, err := cmd.Execute(ctx)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Expected"))
}

func TestCommandEmpty(t *testing.T) {
	ctx := getTestContext("@rocket testcommand")
	cmd := getTestCommand(testHandler)
	_, _, err := cmd.Execute(ctx)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Expected 1 argument(s)"))
}

func TestCommandDuplicateOption(t *testing.T) {
	ctx := getTestContext("@rocket testcommand --myopt=`ayy` --myopt=`letsgo` 1")
	cmd := getTestCommand(testHandler)
	_, _, err := cmd.Execute(ctx)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Duplicate option"))
}

func TestCommandUnrecognizedOption(t *testing.T) {
	ctx := getTestContext("@rocket testcommand --plx=`plox` 1")
	cmd := getTestCommand(testHandler)
	_, _, err := cmd.Execute(ctx)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Unrecognized option"))
}

func TestCommandInvalidArgFormat(t *testing.T) {
	ctx := getTestContext("@rocket testcommand test")
	cmd := getTestCommand(testHandler)
	cmd.Args[0].Format, _ = regexp.Compile("^[0-9]")
	_, _, err := cmd.Execute(ctx)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Invalid format for argument"))
}

func TestCommandTooManyArgs(t *testing.T) {
	ctx := getTestContext("@rocket testcommand test 1 2 3")
	cmd := getTestCommand(testHandler)
	_, _, err := cmd.Execute(ctx)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Too many arguments"))
}

func TestCommandInvalidOptFormat(t *testing.T) {
	ctx := getTestContext("@rocket testcommand --myopt=`test` test")
	cmd := getTestCommand(testHandler)
	cmd.Options["myopt"].Format, _ = regexp.Compile("^[0-9]")
	_, _, err := cmd.Execute(ctx)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Invalid format for option"))
}

func TestCommandHelp(t *testing.T) {
	cmd := getTestCommand(testHandler)
	assert.Equal(t, cmd.Help(), `*Usage:* @rocket testcommand OPTIONS ARGUMENTS

lets go dude!

*Arguments:*
	testarg	what a cool arg

*Options:*
	--myopt	what a sick option
`)
}

func TestCommandMultiWordArg(t *testing.T) {
	ctx := getTestContext("@rocket testcommand lets goo dude!")
	ch := func(context Context) (string, slack.PostMessageParameters) {
		ctx = context
		return context.Message.Text, slack.PostMessageParameters{}
	}
	cmd := getTestCommand(ch)
	cmd.Args[0].MultiWord = true
	res, _, err := cmd.Execute(ctx)
	assert.Nil(t, err)
	assert.Equal(t, res, "@rocket testcommand lets goo dude!")
	assert.Equal(t, ctx.Args[0].Value, "lets goo dude!")
}
