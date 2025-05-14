package gothought

import (
	"context"
	"errors"

	"github.com/gobenpark/gothought/tool"
)

const (
	FinishReasonStop      = "stop"
	FinishReasonToolCalls = "tool_calls"
)

type LanguageModel struct {
	tools         map[string]tool.Tool
	provider      Provider
	messages      []Message
	maxIterations int // maxIterations default int values 10
}

func NewLanguageModel(p Provider, options ...Option) *LanguageModel {
	cli := &LanguageModel{
		provider:      p,
		maxIterations: 10,
		tools:         map[string]tool.Tool{},
	}

	for _, option := range options {
		option(cli)
	}

	return cli
}

func (l *LanguageModel) SetPrompts(prompts []Message) {
	l.messages = prompts
}

// AddTool adds a tool to the client
func (l *LanguageModel) AddTool(t tool.Tool) *LanguageModel {
	l.tools[t.Name()] = t
	return l
}

// SystemPrompt adds a message to the client messages
func (l *LanguageModel) SystemPrompt(prompt string) *LanguageModel {
	l.messages = append(l.messages, Message{
		Role:    "system",
		Message: prompt,
	})
	return l
}

func (l *LanguageModel) AIPrompt(prompt string) *LanguageModel {
	l.messages = append(l.messages, Message{
		Role:    "AI",
		Message: prompt,
	})
	return l
}

func (l *LanguageModel) Prompt(message Message) *LanguageModel {
	l.messages = append(l.messages, message)
	return l
}

func (l *LanguageModel) HumanPrompt(prompt string) *LanguageModel {
	l.messages = append(l.messages, Message{
		Role:    "user",
		Message: prompt,
	})
	return l
}

func (l *LanguageModel) Q(ctx context.Context) (*Message, error) {

	messages := l.messages

	for i := 0; i < l.maxIterations; i++ {
		response, finishReason, err := l.provider.Generate(ctx, l.tools, messages)
		if err != nil {
			return nil, err
		}

		switch finishReason {
		case FinishReasonStop:
			return response, nil
		case FinishReasonToolCalls:
			messages = append(messages, *response)

			for _, tl := range response.ToolCalls {
				tres, err := l.tools[tl.Function.Name].Call(ctx, tl.Function.Arguments)
				if err != nil {
					return nil, err
				}
				messages = append(messages, Message{
					Role:       "tool",
					ToolCallID: tl.ID,
					Message:    tres,
				})
			}
		}
	}
	return nil, errors.New("max iterations reached")
}

func (l *LanguageModel) QStream(ctx context.Context, callback func(Message) error) error {
	if p, ok := any(l.provider).(StreamingCapable); ok {
		return p.GenerateStreaming(ctx, l.messages, callback)
	}

	return errors.New("streaming not supported for this provider")
}

func (o *LanguageModel) QWith(ctx context.Context, oj interface{}) error {
	msgLen := len(o.messages)
	msg := o.messages[msgLen-1]

	msg.Message += "\n\n" + GenerateSchemaPrompt(oj)
	o.messages[msgLen-1] = msg

	res, _, err := o.provider.Generate(ctx, o.tools, o.messages)
	if err != nil {
		return err
	}

	if err := ParsePrompt(oj, res.Message); err != nil {
		return err
	}
	return nil
}
