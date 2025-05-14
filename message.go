package gothought

type Message struct {
	Role       string `json:"role"`
	ToolCallID string `json:"tool_call_id"`
	Message    string
	ToolCalls  []ToolCalls `json:"tool_calls"`
}

type ResponseMessage struct {
	Message string
}
