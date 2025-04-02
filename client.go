package main

type Client struct {
	apikey string
	model  string
}

func NewClient(options ...Option) *Client {
	cli := &Client{}
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
