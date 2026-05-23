package tools

import (
	"os"
	"strings"
)

type listFilesArgs struct {
	Path string `json:"path"`
}

func listFiles(rawArgs string) string {
	var args listFilesArgs
	if err := parseArgs(rawArgs, &args); err != nil {
		return err.Error()
	}
	if args.Path == "" {
		args.Path = "."
	}

	entries, err := os.ReadDir(args.Path)
	if err != nil {
		return "Error listing files: " + err.Error()
	}

	var lines []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}
		lines = append(lines, name)
	}
	return strings.Join(lines, "\n")
}
