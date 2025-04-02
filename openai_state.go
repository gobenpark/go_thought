package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/samber/lo"
)

type OpenAIBody struct {
	Model            string          `json:"model"`
	Messages         []OpenAIMessage `json:"messages"`
	Temperature      float32         `json:"temperature"`
	MaxTokens        int             `json:"max_tokens"`
	TopP             int             `json:"top_p"`
	FrequencyPenalty float32         `json:"frequency_penalty"`
	PresencePenalty  float32         `json:"presence_penalty"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIState struct {
	client   *Client
	messages []Message
}

func NewOpenAIState(client *Client) *OpenAIState {
	return &OpenAIState{client: client, messages: []Message{}}
}

func (o *OpenAIState) SystemPrompt(prompt string) State {
	o.messages = append(o.messages, Message{
		Role:    "system",
		Message: prompt,
	})
	return o
}

func (o *OpenAIState) AIPrompt(prompt string) State {
	o.messages = append(o.messages, Message{
		Role:    "AI",
		Message: prompt,
	})
	return o
}

func (o *OpenAIState) Prompt(message Message) State {
	o.messages = append(o.messages, message)
	return o
}

func (o *OpenAIState) HumanPrompt(prompt string) State {
	o.messages = append(o.messages, Message{
		Role:    "human",
		Message: prompt,
	})
	return o
}

func (o *OpenAIState) Q(ctx context.Context) (string, error) {

	body := OpenAIBody{
		Model: o.client.model,
		Messages: lo.Map(o.messages, func(item Message, index int) OpenAIMessage {
			om := OpenAIMessage{
				Content: item.Message,
				Role:    item.Role,
			}

			switch item.Role {
			case "human":
				om.Role = "user"
			}
			return om
		}),
		Temperature:      0.7,
		MaxTokens:        1024,
		TopP:             1,
		FrequencyPenalty: 0.5,
		PresencePenalty:  0.5,
	}
	bt, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bt))

	request, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(bt))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+o.client.apikey)

	res, err := http.DefaultClient.Do(request.WithContext(ctx))
	if err != nil {
		return "", err
	}
	fmt.Println(res.StatusCode)

	bt, err = io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(bt))
	return string(bt), nil

}
