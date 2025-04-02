package main

type Option func(c *Client)

func WithApiKey(apikey string) Option {
	return func(c *Client) {
		c.apikey = apikey
	}
}

func WithModel(model string) Option {
	return func(c *Client) {
		c.model = model
	}
}
