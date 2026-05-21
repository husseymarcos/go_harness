package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (a *Agent) runTool(call ToolCall) string {
	fmt.Printf("\n[tool] %s %s\n", call.Function.Name, call.Function.Arguments)

	if a.supervision && modifiesSystem(call.Function.Name) {
		fmt.Print("Ejecutar esta acción? [s/n]: ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return "Acción cancelada: no se pudo leer confirmación."
		}
		answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if answer != "s" && answer != "si" && answer != "sí" && answer != "y" {
			return "Acción cancelada por el usuario."
		}
	}

	switch call.Function.Name {
	case "read_file":
		var args struct {
			Path string `json:"path"`
		}
		if err := parseArgs(call, &args); err != nil {
			return err.Error()
		}
		return readFile(args.Path)

	case "write_file":
		var args struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := parseArgs(call, &args); err != nil {
			return err.Error()
		}
		return writeFile(args.Path, args.Content)

	case "run_command":
		var args struct {
			Command string `json:"command"`
		}
		if err := parseArgs(call, &args); err != nil {
			return err.Error()
		}
		return runCommand(args.Command)

	case "list_files":
		var args struct {
			Path string `json:"path"`
		}
		if err := parseArgs(call, &args); err != nil {
			return err.Error()
		}
		return listFiles(args.Path)

	case "web_search":
		var args struct {
			Query string `json:"query"`
		}
		if err := parseArgs(call, &args); err != nil {
			return err.Error()
		}
		return webSearch(args.Query)
	}

	return "Tool desconocida: " + call.Function.Name
}

func parseArgs(call ToolCall, target any) error {
	if strings.TrimSpace(call.Function.Arguments) == "" {
		return errors.New("faltan argumentos")
	}
	return json.Unmarshal([]byte(call.Function.Arguments), target)
}

func modifiesSystem(toolName string) bool {
	return toolName == "write_file" || toolName == "run_command"
}

func readFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return "Error leyendo archivo: " + err.Error()
	}
	return string(content)
}

func writeFile(path, content string) string {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "Error creando directorio: " + err.Error()
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "Error escribiendo archivo: " + err.Error()
	}
	return "Archivo escrito: " + path
}

func runCommand(command string) string {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Comando terminó con error: %v\n%s", err, string(output))
	}
	return string(output)
}

func listFiles(path string) string {
	if path == "" {
		path = "."
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return "Error listando archivos: " + err.Error()
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

func webSearch(query string) string {
	apiKey := os.Getenv("TAVILY_API_KEY")
	if apiKey == "" {
		return "Falta TAVILY_API_KEY. No se puede usar web_search."
	}

	body, _ := json.Marshal(map[string]any{
		"api_key":      apiKey,
		"query":        query,
		"max_results":  5,
		"include_raw":  false,
		"search_depth": "basic",
	})

	res, err := http.Post("https://api.tavily.com/search", "application/json", bytes.NewReader(body))
	if err != nil {
		return "Error buscando en la web: " + err.Error()
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return "Error leyendo respuesta de Tavily: " + err.Error()
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Sprintf("Tavily devolvió status %d: %s", res.StatusCode, string(raw))
	}
	return string(raw)
}
