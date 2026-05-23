package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const systemPrompt = `You are a minimalist coding agent.
Your job is to help the user using tools when needed.
Before writing files or running commands, briefly explain what you are doing.
When you finish a task, respond without requesting tools.`

type Agent struct {
	apiKey      string
	model       string
	client      *http.Client
	messages    []Message
	planMode    bool
	supervision bool
}

func NewAgent(apiKey, model string, client *http.Client) *Agent {
	return &Agent{
		apiKey: apiKey,
		model:  model,
		client: client,
		messages: []Message{
			{Role: "system", Content: systemPrompt},
		},
	}
}

func (a *Agent) HandleCommand(input string) bool {
	switch input {
	case ":quit", ":exit":
		os.Exit(0)
	case ":plan on":
		a.planMode = true
		fmt.Println("Plan mode activated.")
		return true
	case ":plan off":
		a.planMode = false
		fmt.Println("Plan mode deactivated.")
		return true
	case ":supervision on":
		a.supervision = true
		fmt.Println("Supervision activated.")
		return true
	case ":supervision off":
		a.supervision = false
		fmt.Println("Supervision deactivated.")
		return true
	}
	return false
}

func (a *Agent) RunUserTurn(input string) error {
	if a.planMode {
		ok, replacement, err := a.confirmPlan(input)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("Task cancelled.")
			return nil
		}
		if replacement != "" {
			input = replacement
		}
	}

	a.messages = append(a.messages, Message{Role: "user", Content: input})

	for iteration := 1; ; iteration++ {
		response, err := a.chat(a.messages, tools())
		if err != nil {
			return err
		}
		if len(response.Choices) == 0 {
			return errors.New("the model returned no responses")
		}

		message := response.Choices[0].Message
		a.messages = append(a.messages, message)

		if len(message.ToolCalls) == 0 {
			fmt.Printf("\n%s\n", message.Content)
			fmt.Printf("(internal loop iterations: %d)\n", iteration)
			return nil
		}

		for _, call := range message.ToolCalls {
			result := a.runTool(call)
			a.messages = append(a.messages, Message{
				Role:       "tool",
				ToolCallID: call.ID,
				ToolName:   call.Function.Name,
				Content:    result,
			})
		}
	}
}

func (a *Agent) confirmPlan(input string) (bool, string, error) {
	planMessages := []Message{
		{Role: "system", Content: "Make a brief, numbered, and concrete plan. Do not use tools."},
		{Role: "user", Content: input},
	}
	response, err := a.chat(planMessages, nil)
	if err != nil {
		return false, "", err
	}
	if len(response.Choices) == 0 {
		return false, "", errors.New("the model returned no plan")
	}

	fmt.Println("\nProposed plan:")
	fmt.Println(response.Choices[0].Message.Content)
	fmt.Print("\nApprove? [y/n/modify]: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false, "", scanner.Err()
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))

	if answer == "y" || answer == "yes" {
		return true, "", nil
	}
	if answer == "modify" || answer == "m" {
		fmt.Print("New instruction: ")
		if !scanner.Scan() {
			return false, "", scanner.Err()
		}
		return true, strings.TrimSpace(scanner.Text()), nil
	}
	return false, "", nil
}
