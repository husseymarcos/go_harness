package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Falta OPENAI_API_KEY en el entorno.")
		os.Exit(1)
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4.1-mini"
	}

	agent := NewAgent(apiKey, model, &http.Client{})

	fmt.Println("Coding agent listo.")
	fmt.Println("Comandos: :plan on/off, :supervision on/off, :quit")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if agent.HandleCommand(input) {
			continue
		}

		if err := agent.RunUserTurn(input); err != nil {
			fmt.Println("Error:", err)
		}
	}
}
