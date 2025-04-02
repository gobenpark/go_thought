package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/samber/lo"
)

type ClaudeInput struct {
	Model    string `json:"model"`
	MaxToken int
	Messages []Message
}

type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"content"`
	Id           string      `json:"id"`
	Model        string      `json:"model"`
	Role         string      `json:"role"`
	StopReason   string      `json:"stop_reason"`
	StopSequence interface{} `json:"stop_sequence"`
	Type         string      `json:"type"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

type ClaudeState struct {
	client   *Client
	messages []Message
}

func NewClaudeState(client *Client) *ClaudeState {
	return &ClaudeState{client: client, messages: []Message{}}
}

func (c *ClaudeState) SystemPrompt(prompt string) State {
	c.messages = append(c.messages, Message{
		Role:    "system",
		Message: prompt,
	})
	return c
}

func (c *ClaudeState) Prompt(message Message) State {
	c.messages = append(c.messages, message)
	return c
}

func (c *ClaudeState) HumanPrompt(prompt string) State {
	c.messages = append(c.messages, Message{
		Role:    "user",
		Message: prompt,
	})
	return c
}

func (c *ClaudeState) AIPrompt(prompt string) State {
	c.messages = append(c.messages, Message{
		Role:    "assistant",
		Message: prompt,
	})
	return c
}

func (c *ClaudeState) Q(ctx context.Context) ([]ResponseMessage, error) {
	body := OpenAIBody{
		Model: c.client.model,
		Messages: lo.Map(c.messages, func(item Message, index int) OpenAIMessage {
			return OpenAIMessage{
				Content: item.Message,
				Role:    item.Role,
			}
		}),
		MaxTokens: 1024,
	}
	bt, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(bt))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", c.client.apikey)
	request.Header.Set("anthropic-version", "2023-06-01")

	res, err := http.DefaultClient.Do(request.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		buf := bytes.Buffer{}
		io.Copy(&buf, res.Body)
		return nil, errors.New(buf.String())
	}

	var response ClaudeResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	var messages []ResponseMessage
	for _, i := range response.Content {
		messages = append(messages, ResponseMessage{
			Message: i.Text,
		})
	}

	return messages, nil
}

func (c *ClaudeState) QStream(ctx context.Context, callback func(ResponseMessage) error) error {
	//TODO implement me
	panic("implement me")
}
