package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	ToolName   string     `json:"tool_name,omitempty"`
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

type GeminiRequest struct {
	SystemInstruction *GeminiContent  `json:"systemInstruction,omitempty"`
	Contents          []GeminiContent `json:"contents"`
	Tools             []GeminiTool    `json:"tools,omitempty"`
	GenerationConfig  GeminiConfig    `json:"generationConfig"`
}

type GeminiConfig struct {
	Temperature float64 `json:"temperature"`
}

type GeminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text             string                  `json:"text,omitempty"`
	FunctionCall     *GeminiFunctionCall     `json:"functionCall,omitempty"`
	FunctionResponse *GeminiFunctionResponse `json:"functionResponse,omitempty"`
}

type GeminiFunctionCall struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args,omitempty"`
}

type GeminiFunctionResponse struct {
	Name     string         `json:"name"`
	Response map[string]any `json:"response"`
}

type GeminiTool struct {
	FunctionDeclarations []ToolSchema `json:"functionDeclarations"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content GeminiContent `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func (a *Agent) chat(messages []Message, toolDefs []ToolDef) (ChatResponse, error) {
	body := GeminiRequest{
		SystemInstruction: systemInstruction(messages),
		Contents:          geminiContents(messages),
		Tools:             geminiTools(toolDefs),
		GenerationConfig: GeminiConfig{
			Temperature: 0,
		},
	}

	rawBody, err := json.Marshal(body)
	if err != nil {
		return ChatResponse{}, err
	}

	endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", url.PathEscape(a.model), url.QueryEscape(a.apiKey))
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(rawBody))
	if err != nil {
		return ChatResponse{}, err
	}
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

	var decoded GeminiResponse
	if err := json.Unmarshal(rawResponse, &decoded); err != nil {
		return ChatResponse{}, err
	}
	if decoded.Error != nil {
		return ChatResponse{}, errors.New(decoded.Error.Message)
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return ChatResponse{}, fmt.Errorf("status %d: %s", res.StatusCode, string(rawResponse))
	}
	return chatResponse(decoded), nil
}

func systemInstruction(messages []Message) *GeminiContent {
	for _, message := range messages {
		if message.Role == "system" && message.Content != "" {
			return &GeminiContent{Parts: []GeminiPart{{Text: message.Content}}}
		}
	}
	return nil
}

func geminiContents(messages []Message) []GeminiContent {
	contents := make([]GeminiContent, 0, len(messages))
	for _, message := range messages {
		switch message.Role {
		case "system":
			continue
		case "assistant":
			contents = append(contents, GeminiContent{Role: "model", Parts: geminiMessageParts(message)})
		case "tool":
			contents = append(contents, GeminiContent{Role: "function", Parts: []GeminiPart{{FunctionResponse: &GeminiFunctionResponse{
				Name: message.ToolName,
				Response: map[string]any{
					"result": message.Content,
				},
			}}}})
		default:
			contents = append(contents, GeminiContent{Role: "user", Parts: []GeminiPart{{Text: message.Content}}})
		}
	}
	return contents
}

func geminiMessageParts(message Message) []GeminiPart {
	parts := make([]GeminiPart, 0, 1+len(message.ToolCalls))
	if message.Content != "" {
		parts = append(parts, GeminiPart{Text: message.Content})
	}
	for _, call := range message.ToolCalls {
		parts = append(parts, GeminiPart{FunctionCall: &GeminiFunctionCall{
			Name: call.Function.Name,
			Args: functionArgs(call.Function.Arguments),
		}})
	}
	return parts
}

func functionArgs(raw string) map[string]any {
	var args map[string]any
	if err := json.Unmarshal([]byte(raw), &args); err != nil {
		return map[string]any{}
	}
	return args
}

func geminiTools(toolDefs []ToolDef) []GeminiTool {
	if len(toolDefs) == 0 {
		return nil
	}
	declarations := make([]ToolSchema, 0, len(toolDefs))
	for _, toolDef := range toolDefs {
		declarations = append(declarations, toolDef.Function)
	}
	return []GeminiTool{{FunctionDeclarations: declarations}}
}

func chatResponse(response GeminiResponse) ChatResponse {
	var converted ChatResponse
	if len(response.Candidates) == 0 {
		return converted
	}

	message := Message{Role: "assistant"}
	for i, part := range response.Candidates[0].Content.Parts {
		if part.Text != "" {
			message.Content += part.Text
		}
		if part.FunctionCall != nil {
			rawArgs, _ := json.Marshal(part.FunctionCall.Args)
			message.ToolCalls = append(message.ToolCalls, ToolCall{
				ID:   fmt.Sprintf("gemini-call-%d", i),
				Type: "function",
				Function: FunctionCall{
					Name:      part.FunctionCall.Name,
					Arguments: string(rawArgs),
				},
			})
		}
	}
	converted.Choices = append(converted.Choices, struct {
		Message Message `json:"message"`
	}{Message: message})
	return converted
}
