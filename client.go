package main

import (
	"github.com/gobenpark/gothought/tool"
)

type LanguageModel struct {
	apikey string
	model  string
	tools  []tool.Tool
}

func NewLanguageModel(options ...Option) *LanguageModel {
	cli := &LanguageModel{
		tools: []tool.Tool{},
	}
	for _, option := range options {
		option(cli)
	}

	return cli
}

func (c *LanguageModel) State(st StateType) State {
	switch st {
	case OPENAI:
		return NewOpenAIState(c)
	case CLAUDE:
		return NewClaudeState(c)
	}
	return nil
}

// AddTool adds a tool to the client
func (c *LanguageModel) AddTool(t tool.Tool) *LanguageModel {
	c.tools = append(c.tools, t)
	return c
}
