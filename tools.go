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

func tools() []ToolSchema {
	return []ToolSchema{
		makeTool("read_file", "Read the content of a file.", map[string]any{
			"path": stringParam("Path of the file to read."),
		}, []string{"path"}),
		makeTool("write_file", "Write content to a file, replacing the current content.", map[string]any{
			"path":    stringParam("Path of the file to write."),
			"content": stringParam("Full content of the file."),
		}, []string{"path", "content"}),
		makeTool("run_command", "Run a terminal command and return stdout and stderr.", map[string]any{
			"command": stringParam("Full command to run."),
		}, []string{"command"}),
		makeTool("list_files", "List the files in a directory.", map[string]any{
			"path": stringParam("Directory to list. Use . if not specified."),
		}, []string{"path"}),
		makeTool("web_search", "Search the web with Tavily.", map[string]any{
			"query": stringParam("Search query."),
		}, []string{"query"}),
	}
}

func makeTool(name, description string, properties map[string]any, required []string) ToolSchema {
	return ToolSchema{
		Name:        name,
		Description: description,
		Parameters: map[string]any{
			"type":       "object",
			"properties": properties,
			"required":   required,
		},
	}
}

func stringParam(description string) map[string]string {
	return map[string]string{
		"type":        "string",
		"description": description,
	}
}
