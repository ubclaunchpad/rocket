package bot

type Responder func(msg string) string

type Conversation struct {
	entry Interaction
}

type Interaction struct {
}
