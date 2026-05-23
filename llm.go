package main

type Message struct {
	Role      string
	Content   string
	ToolCalls []ToolCall
	ToolName  string
}

type ToolCall struct {
	Function FunctionCall
}

type FunctionCall struct {
	Name      string
	Arguments string
}

type ToolSchema struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type Provider interface {
	Chat(messages []Message, tools []ToolSchema) (*Message, error)
}
