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

type GeminiProvider struct {
	APIKey string
	Model  string
	Client *http.Client
}

type geminiRequest struct {
	SystemInstruction *geminiContent  `json:"systemInstruction,omitempty"`
	Contents          []geminiContent `json:"contents"`
	Tools             []geminiTool    `json:"tools,omitempty"`
	GenerationConfig  geminiConfig    `json:"generationConfig"`
}

type geminiConfig struct {
	Temperature float64 `json:"temperature"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text             string                  `json:"text,omitempty"`
	FunctionCall     *geminiFunctionCall     `json:"functionCall,omitempty"`
	FunctionResponse *geminiFunctionResponse `json:"functionResponse,omitempty"`
}

type geminiFunctionCall struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args,omitempty"`
}

type geminiFunctionResponse struct {
	Name     string         `json:"name"`
	Response map[string]any `json:"response"`
}

type geminiTool struct {
	FunctionDeclarations []ToolSchema `json:"functionDeclarations"`
}

type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (p *GeminiProvider) Chat(messages []Message, toolDefs []ToolSchema) (*Message, error) {
	contents, err := geminiContents(messages)
	if err != nil {
		return nil, err
	}

	body := geminiRequest{
		SystemInstruction: systemInstruction(messages),
		Contents:          contents,
		Tools:             geminiTools(toolDefs),
		GenerationConfig: geminiConfig{
			Temperature: 0,
		},
	}

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", url.PathEscape(p.Model), url.QueryEscape(p.APIKey))
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(rawBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	rawResponse, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var decoded geminiResponse
	if err := json.Unmarshal(rawResponse, &decoded); err != nil {
		return nil, err
	}
	if decoded.Error != nil {
		return nil, errors.New(decoded.Error.Message)
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("status %d: %s", res.StatusCode, string(rawResponse))
	}
	return messageFromGemini(decoded)
}

func systemInstruction(messages []Message) *geminiContent {
	for _, message := range messages {
		if message.Role == "system" && message.Content != "" {
			return &geminiContent{Parts: []geminiPart{{Text: message.Content}}}
		}
	}
	return nil
}

func geminiContents(messages []Message) ([]geminiContent, error) {
	contents := make([]geminiContent, 0, len(messages))
	for _, message := range messages {
		switch message.Role {
		case "system":
			continue
		case "assistant":
			parts, err := geminiMessageParts(message)
			if err != nil {
				return nil, err
			}
			contents = append(contents, geminiContent{Role: "model", Parts: parts})
		case "tool":
			contents = append(contents, geminiContent{Role: "function", Parts: []geminiPart{{FunctionResponse: &geminiFunctionResponse{
				Name: message.ToolName,
				Response: map[string]any{
					"result": message.Content,
				},
			}}}})
		default:
			contents = append(contents, geminiContent{Role: "user", Parts: []geminiPart{{Text: message.Content}}})
		}
	}
	return contents, nil
}

func geminiMessageParts(message Message) ([]geminiPart, error) {
	parts := make([]geminiPart, 0, 1+len(message.ToolCalls))
	if message.Content != "" {
		parts = append(parts, geminiPart{Text: message.Content})
	}
	for _, call := range message.ToolCalls {
		args, err := functionArgs(call.Function.Arguments)
		if err != nil {
			return nil, err
		}
		parts = append(parts, geminiPart{FunctionCall: &geminiFunctionCall{
			Name: call.Function.Name,
			Args: args,
		}})
	}
	return parts, nil
}

func functionArgs(raw string) (map[string]any, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(raw), &args); err != nil {
		return nil, err
	}
	return args, nil
}

func geminiTools(toolDefs []ToolSchema) []geminiTool {
	if len(toolDefs) == 0 {
		return nil
	}
	return []geminiTool{{FunctionDeclarations: toolDefs}}
}

func messageFromGemini(response geminiResponse) (*Message, error) {
	if len(response.Candidates) == 0 {
		return nil, nil
	}

	message := Message{Role: "assistant"}
	for _, part := range response.Candidates[0].Content.Parts {
		if part.Text != "" {
			message.Content += part.Text
		}
		if part.FunctionCall != nil {
			rawArgs, err := json.Marshal(part.FunctionCall.Args)
			if err != nil {
				return nil, err
			}
			message.ToolCalls = append(message.ToolCalls, ToolCall{
				Function: FunctionCall{
					Name:      part.FunctionCall.Name,
					Arguments: string(rawArgs),
				},
			})
		}
	}
	return &message, nil
}
