package tools

import (
	"os"
	"path/filepath"
)

type writeFileArgs struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func writeFile(rawArgs string) string {
	var args writeFileArgs
	if err := parseArgs(rawArgs, &args); err != nil {
		return err.Error()
	}

	if err := os.MkdirAll(filepath.Dir(args.Path), 0755); err != nil {
		return "Error creating directory: " + err.Error()
	}
	if err := os.WriteFile(args.Path, []byte(args.Content), 0644); err != nil {
		return "Error writing file: " + err.Error()
	}
	return "File written: " + args.Path
}
