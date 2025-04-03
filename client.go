package main

import (
	"github.com/gobenpark/gothought/tool"
)

type Client struct {
	apikey string
	model  string
	tools  []tool.Tool
}

func NewClient(options ...Option) *Client {
	cli := &Client{
		tools: []tool.Tool{},
	}
	for _, option := range options {
		option(cli)
	}

	return cli
}

func (c *Client) State(st StateType) State {
	switch st {
	case OPENAI:
		return NewOpenAIState(c)
	case CLAUDE:
		return NewClaudeState(c)
	}
	return nil
}

// AddTool adds a tool to the client
func (c *Client) AddTool(t tool.Tool) *Client {
	c.tools = append(c.tools, t)
	return c
}
