package tools

import "os"

type readFileArgs struct {
	Path string `json:"path"`
}

func readFile(rawArgs string) string {
	var args readFileArgs
	if err := parseArgs(rawArgs, &args); err != nil {
		return err.Error()
	}

	content, err := os.ReadFile(args.Path)
	if err != nil {
		return "Error reading file: " + err.Error()
	}
	return string(content)
}
