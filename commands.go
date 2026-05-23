package main

import (
	"fmt"
	"os"
)

func (a *Agent) HandleCommand(input string) bool {
	switch input {
	case ":quit", ":exit":
		os.Exit(0)
	case ":plan on":
		a.planMode = true
		fmt.Println("Plan mode activated.")
		return true
	case ":plan off":
		a.planMode = false
		fmt.Println("Plan mode deactivated.")
		return true
	case ":supervision on":
		a.supervision = true
		fmt.Println("Supervision activated.")
		return true
	case ":supervision off":
		a.supervision = false
		fmt.Println("Supervision deactivated.")
		return true
	}
	return false
}
