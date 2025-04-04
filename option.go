package main

type Option func(c *LanguageModel)

func WithApiKey(apikey string) Option {
	return func(c *LanguageModel) {
		c.apikey = apikey
	}
}

func WithModel(model string) Option {
	return func(c *LanguageModel) {
		c.model = model
	}
}
