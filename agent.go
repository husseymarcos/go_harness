package main

import (
	"errors"
	"fmt"
)

const systemPrompt = `You are a minimalist coding agent.
Your job is to help the user using tools when needed.
Before writing files or running commands, briefly explain what you are doing.
When you finish a task, respond without requesting tools.`

type Agent struct {
	provider    Provider
	messages    []Message
	planMode    bool
	supervision bool
}

func NewAgent(provider Provider) *Agent {
	return &Agent{
		provider: provider,
		messages: []Message{
			{Role: "system", Content: systemPrompt},
		},
	}
}

func (a *Agent) RunUserTurn(input string) error {
	input, ok, err := a.prepareInput(input)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println("Task cancelled.")
		return nil
	}

	a.addMessage("user", input)

	for iteration := 1; ; iteration++ {
		message, err := a.askModel()
		if err != nil {
			return err
		}

		if message.hasFinalAnswer() {
			printFinalAnswer(message.Content, iteration)
			return nil
		}

		a.runToolCalls(message.ToolCalls)
	}
}

func (a *Agent) addMessage(role, content string) {
	a.messages = append(a.messages, Message{Role: role, Content: content})
}

func (a *Agent) askModel() (*Message, error) {
	message, err := a.provider.Chat(a.messages, tools())
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, errors.New("the model returned no responses")
	}

	a.messages = append(a.messages, *message)
	return message, nil
}

func (m *Message) hasFinalAnswer() bool {
	return len(m.ToolCalls) == 0
}

func printFinalAnswer(content string, iterations int) {
	fmt.Printf("\n%s\n", content)
	fmt.Printf("(internal loop iterations: %d)\n", iterations)
}

func (a *Agent) runToolCalls(calls []ToolCall) {
	for _, call := range calls {
		result := a.runTool(call)
		a.messages = append(a.messages, Message{
			Role:     "tool",
			ToolName: call.Function.Name,
			Content:  result,
		})
	}
}
