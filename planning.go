package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func (a *Agent) prepareInput(input string) (string, bool, error) {
	if !a.planMode {
		return input, true, nil
	}

	ok, replacement, err := a.confirmPlan(input)
	if err != nil || !ok {
		return input, ok, err
	}
	if replacement != "" {
		return replacement, true, nil
	}
	return input, true, nil
}

func (a *Agent) confirmPlan(input string) (bool, string, error) {
	planMessages := []Message{
		{Role: "system", Content: "Make a brief, numbered, and concrete plan. Do not use tools."},
		{Role: "user", Content: input},
	}
	message, err := a.provider.Chat(planMessages, nil)
	if err != nil {
		return false, "", err
	}
	if message == nil {
		return false, "", errors.New("the model returned no plan")
	}

	fmt.Println("\nProposed plan:")
	fmt.Println(message.Content)
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
