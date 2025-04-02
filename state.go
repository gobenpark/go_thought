package main

import "context"

type StateType string

const (
	OPENAI StateType = "openai"
)

type State interface {
	SystemPrompt(prompt string) State
	Prompt(message Message) State
	HumanPrompt(prompt string) State
	AIPrompt(prompt string) State
	Q(ctx context.Context) ([]Message, error)
	QStream(ctx context.Context, callback func(Message) error) error
}
