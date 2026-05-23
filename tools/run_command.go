package tools

import (
	"fmt"
	"os/exec"
)

type runCommandArgs struct {
	Command string `json:"command"`
}

func runCommand(rawArgs string) string {
	var args runCommandArgs
	if err := parseArgs(rawArgs, &args); err != nil {
		return err.Error()
	}

	cmd := exec.Command("sh", "-c", args.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Command failed with error: %v\n%s", err, string(output))
	}
	return string(output)
}
