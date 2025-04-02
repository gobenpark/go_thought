package main

import (
	"bufio"
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
	Stream           bool            `json:"stream"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role        string        `json:"role"`
			Content     string        `json:"content"`
			Refusal     interface{}   `json:"refusal"`
			Annotations []interface{} `json:"annotations"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens        int `json:"prompt_tokens"`
		CompletionTokens    int `json:"completion_tokens"`
		TotalTokens         int `json:"total_tokens"`
		PromptTokensDetails struct {
			CachedTokens int `json:"cached_tokens"`
			AudioTokens  int `json:"audio_tokens"`
		} `json:"prompt_tokens_details"`
		CompletionTokensDetails struct {
			ReasoningTokens          int `json:"reasoning_tokens"`
			AudioTokens              int `json:"audio_tokens"`
			AcceptedPredictionTokens int `json:"accepted_prediction_tokens"`
			RejectedPredictionTokens int `json:"rejected_prediction_tokens"`
		} `json:"completion_tokens_details"`
	} `json:"usage"`
	ServiceTier       string `json:"service_tier"`
	SystemFingerprint string `json:"system_fingerprint"`
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

func (o *OpenAIState) Q(ctx context.Context) ([]Message, error) {

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
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(bt))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+o.client.apikey)

	res, err := http.DefaultClient.Do(request.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var result OpenAIResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var responseMessage []Message
	for _, i := range result.Choices {
		responseMessage = append(responseMessage, Message{
			Role:    i.Message.Role,
			Message: i.Message.Content,
		})
	}

	return responseMessage, nil
}

func (o *OpenAIState) Q_Stream(ctx context.Context, callback func(Message) error) error {
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
		Stream:           true,
	}

	bt, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(bt))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+o.client.apikey)

	res, err := http.DefaultClient.Do(request.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("API returned non-200 status code: %d, body: %s", res.StatusCode, string(bodyBytes))
	}

	reader := bufio.NewReader(res.Body)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, []byte("data: ")) {
			continue
		}

		data := line[6:]

		if string(data) == "[DONE]" {
			break
		}

		var chunkResponse struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
					Role    string `json:"role"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.Unmarshal(data, &chunkResponse); err != nil {
			return err
		}

		if len(chunkResponse.Choices) > 0 && chunkResponse.Choices[0].Delta.Content != "" {
			message := Message{
				Role:    "assistant",
				Message: chunkResponse.Choices[0].Delta.Content,
			}

			if err := callback(message); err != nil {
				return err
			}
		}
	}

	return nil
}
