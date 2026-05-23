package tools

import (
	"encoding/json"
	"errors"
	"strings"
)

type handler func(string) string

var handlers = map[string]handler{
	"read_file":   readFile,
	"write_file":  writeFile,
	"run_command": runCommand,
	"list_files":  listFiles,
	"web_search":  webSearch,
}

func Run(name, arguments string) string {
	handler, ok := handlers[name]
	if !ok {
		return "Unknown tool: " + name
	}
	return handler(arguments)
}

func ModifiesSystem(name string) bool {
	return name == "write_file" || name == "run_command"
}

func parseArgs(raw string, target any) error {
	if strings.TrimSpace(raw) == "" {
		return errors.New("missing arguments")
	}
	return json.Unmarshal([]byte(raw), target)
}
