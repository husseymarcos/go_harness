package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Tools       []ToolDef `json:"tools,omitempty"`
	ToolChoice  string    `json:"tool_choice,omitempty"`
	Temperature float64   `json:"temperature"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (a *Agent) chat(messages []Message, toolDefs []ToolDef) (ChatResponse, error) {
	body := ChatRequest{
		Model:       a.model,
		Messages:    messages,
		Tools:       toolDefs,
		Temperature: 0,
	}
	if len(toolDefs) > 0 {
		body.ToolChoice = "auto"
	}

	rawBody, err := json.Marshal(body)
	if err != nil {
		return ChatResponse{}, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(rawBody))
	if err != nil {
		return ChatResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := a.client.Do(req)
	if err != nil {
		return ChatResponse{}, err
	}
	defer res.Body.Close()

	rawResponse, err := io.ReadAll(res.Body)
	if err != nil {
		return ChatResponse{}, err
	}

	var decoded ChatResponse
	if err := json.Unmarshal(rawResponse, &decoded); err != nil {
		return ChatResponse{}, err
	}
	if decoded.Error != nil {
		return ChatResponse{}, errors.New(decoded.Error.Message)
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return ChatResponse{}, fmt.Errorf("status %d: %s", res.StatusCode, string(rawResponse))
	}
	return decoded, nil
}
