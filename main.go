package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const geminiModel = "gemini-2.0-flash"

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env:", err)
		os.Exit(1)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("Missing GEMINI_API_KEY in the environment.")
		os.Exit(1)
	}

	provider := &GeminiProvider{
		APIKey: apiKey,
		Model:  geminiModel,
		Client: &http.Client{},
	}
	agent := NewAgent(provider)

	fmt.Println("Coding agent ready.")
	fmt.Println("Commands: :plan on/off, :supervision on/off, :quit")

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
