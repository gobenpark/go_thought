package gothought

/*
	Claude Input Model

API Spec

	{
		"model": "claude-3-7-sonnet-20250219",
		"max_tokens": 1024,
		"messages": [
			{"role": "user", "content": "Hello, world"}
		]
	}
*/
//type ClaudeInput struct {
//	Model    string `json:"model"`
//	MaxToken int
//	Messages []Message
//}
//
//type ClaudeResponse struct {
//	Content []struct {
//		Text string `json:"text"`
//		Type string `json:"type"`
//	} `json:"content"`
//	Id           string      `json:"id"`
//	Model        string      `json:"model"`
//	Role         string      `json:"role"`
//	StopReason   string      `json:"stop_reason"`
//	StopSequence interface{} `json:"stop_sequence"`
//	Type         string      `json:"type"`
//	Usage        struct {
//		InputTokens  int `json:"input_tokens"`
//		OutputTokens int `json:"output_tokens"`
//	} `json:"usage"`
//}
//
//type ClaudeState struct {
//	client   *LanguageModel
//	messages []Message
//}
//
//func NewClaudeState(client *LanguageModel) *ClaudeState {
//	return &ClaudeState{client: client, messages: []Message{}}
//}
//
//func (c *ClaudeState) SystemPrompt(prompt string) State {
//	c.messages = append(c.messages, Message{
//		Role:    "system",
//		Message: prompt,
//	})
//	return c
//}
//
//func (c *ClaudeState) Prompt(message Message) State {
//	c.messages = append(c.messages, message)
//	return c
//}
//
//func (c *ClaudeState) HumanPrompt(prompt string) State {
//	c.messages = append(c.messages, Message{
//		Role:    "user",
//		Message: prompt,
//	})
//	return c
//}
//
//func (c *ClaudeState) AIPrompt(prompt string) State {
//	c.messages = append(c.messages, Message{
//		Role:    "assistant",
//		Message: prompt,
//	})
//	return c
//}
//
//func (c *ClaudeState) QWithType(ctx context.Context, oj interface{}) error {
//	return nil
//}
//
//func (c *ClaudeState) Q(ctx context.Context) ([]ResponseMessage, error) {
//	body := ClaudeInput{
//		Model:    c.client.model,
//		Messages: c.messages,
//		MaxToken: 1024,
//	}
//	bt, err := json.Marshal(body)
//	if err != nil {
//		return nil, err
//	}
//
//	request, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(bt))
//	if err != nil {
//		return nil, err
//	}
//
//	request.Header.Set("Content-Type", "application/json")
//	request.Header.Set("x-api-key", c.client.apikey)
//	request.Header.Set("anthropic-version", "2023-06-01")
//
//	res, err := http.DefaultClient.Do(request.WithContext(ctx))
//	if err != nil {
//		return nil, err
//	}
//
//	if res.StatusCode != 200 {
//		buf := bytes.Buffer{}
//		if _, err := io.Copy(&buf, res.Body); err != nil {
//			return nil, err
//		}
//		return nil, errors.New(buf.String())
//	}
//
//	var response ClaudeResponse
//	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
//		return nil, err
//	}
//
//	var messages []ResponseMessage
//	for _, i := range response.Content {
//		messages = append(messages, ResponseMessage{
//			Message: i.Text,
//		})
//	}
//
//	return nil, nil
//}
//
//func (c *ClaudeState) QStream(ctx context.Context, callback func(ResponseMessage) error) error {
//	//TODO implement me
//	panic("implement me")
//}
