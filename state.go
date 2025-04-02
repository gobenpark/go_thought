package main

import "context"

type StateType string

const (
	OPENAI StateType = "openai"
	CLAUDE StateType = "claude"
)

type State interface {
	SystemPrompt(prompt string) State
	Prompt(message Message) State
	HumanPrompt(prompt string) State
	AIPrompt(prompt string) State
	Q(ctx context.Context) ([]ResponseMessage, error)
	QStream(ctx context.Context, callback func(ResponseMessage) error) error
}
