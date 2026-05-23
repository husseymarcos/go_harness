package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	toolimpl "husseymarcos/go_harness/tools"
)

func (a *Agent) runTool(call ToolCall) string {
	fmt.Printf("\n[tool] %s %s\n", call.Function.Name, call.Function.Arguments)

	if a.supervision && toolimpl.ModifiesSystem(call.Function.Name) {
		fmt.Print("Execute this action? [y/n]: ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return "Action cancelled: could not read confirmation."
		}
		answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if answer != "y" && answer != "yes" {
			return "Action cancelled by the user."
		}
	}

	return toolimpl.Run(call.Function.Name, call.Function.Arguments)
}
