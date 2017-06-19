package main

import "github.com/ubclaunchpad/rocket/model"
import "github.com/nlopes/slack"

const (
	stateHello = iota
	stateName
	stateEmail
	stateProgram
)

// Session manages a conversation between a Slack user and Rocket
// across multiple messages.
type Session struct {
	member  model.Member
	state   int
	msgChan chan *slack.Message
}

func (s *Session) Start(initMsg *slack.MessageEvent) {
	s.state = stateHello
}
