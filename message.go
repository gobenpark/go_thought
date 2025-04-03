package main

type Message struct {
	Role       string `json:"role"`
	ToolCallID string `json:"tool_call_id"`
	Message    string
}

type ResponseMessage struct {
	Message string
}
