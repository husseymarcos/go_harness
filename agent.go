package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const systemPrompt = `Sos un coding agent minimalista.
Tu trabajo es ayudar al usuario usando herramientas cuando haga falta.
Antes de escribir archivos o ejecutar comandos, explicá brevemente qué estás haciendo.
Cuando termines una tarea, respondé sin pedir herramientas.`

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
		fmt.Println("Plan mode activado.")
		return true
	case ":plan off":
		a.planMode = false
		fmt.Println("Plan mode desactivado.")
		return true
	case ":supervision on":
		a.supervision = true
		fmt.Println("Supervisión activada.")
		return true
	case ":supervision off":
		a.supervision = false
		fmt.Println("Supervisión desactivada.")
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
			fmt.Println("Tarea cancelada.")
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
			return errors.New("el modelo no devolvió respuestas")
		}

		message := response.Choices[0].Message
		a.messages = append(a.messages, message)

		if len(message.ToolCalls) == 0 {
			fmt.Printf("\n%s\n", message.Content)
			fmt.Printf("(iteraciones del loop interno: %d)\n", iteration)
			return nil
		}

		for _, call := range message.ToolCalls {
			result := a.runTool(call)
			a.messages = append(a.messages, Message{
				Role:       "tool",
				ToolCallID: call.ID,
				Content:    result,
			})
		}
	}
}

func (a *Agent) confirmPlan(input string) (bool, string, error) {
	planMessages := []Message{
		{Role: "system", Content: "Armá un plan breve, numerado y concreto. No uses tools."},
		{Role: "user", Content: input},
	}
	response, err := a.chat(planMessages, nil)
	if err != nil {
		return false, "", err
	}
	if len(response.Choices) == 0 {
		return false, "", errors.New("el modelo no devolvió plan")
	}

	fmt.Println("\nPlan propuesto:")
	fmt.Println(response.Choices[0].Message.Content)
	fmt.Print("\nAprobar? [s/n/modificar]: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false, "", scanner.Err()
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))

	if answer == "s" || answer == "si" || answer == "sí" || answer == "y" {
		return true, "", nil
	}
	if answer == "modificar" || answer == "m" {
		fmt.Print("Nueva instrucción: ")
		if !scanner.Scan() {
			return false, "", scanner.Err()
		}
		return true, strings.TrimSpace(scanner.Text()), nil
	}
	return false, "", nil
}
