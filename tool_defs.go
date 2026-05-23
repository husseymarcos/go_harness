package main

type ToolDef struct {
	Type     string     `json:"type"`
	Function ToolSchema `json:"function"`
}

type ToolSchema struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

func tools() []ToolDef {
	return []ToolDef{
		tool("read_file", "Read the content of a file.", map[string]any{
			"path": stringParam("Path of the file to read."),
		}, []string{"path"}),
		tool("write_file", "Write content to a file, replacing the current content.", map[string]any{
			"path":    stringParam("Path of the file to write."),
			"content": stringParam("Full content of the file."),
		}, []string{"path", "content"}),
		tool("run_command", "Run a terminal command and return stdout and stderr.", map[string]any{
			"command": stringParam("Full command to run."),
		}, []string{"command"}),
		tool("list_files", "List the files in a directory.", map[string]any{
			"path": stringParam("Directory to list. Use . if not specified."),
		}, []string{"path"}),
		tool("web_search", "Search the web with Tavily.", map[string]any{
			"query": stringParam("Search query."),
		}, []string{"query"}),
	}
}

func tool(name, description string, properties map[string]any, required []string) ToolDef {
	return ToolDef{
		Type: "function",
		Function: ToolSchema{
			Name:        name,
			Description: description,
			Parameters: map[string]any{
				"type":       "object",
				"properties": properties,
				"required":   required,
			},
		},
	}
}

func stringParam(description string) map[string]string {
	return map[string]string{
		"type":        "string",
		"description": description,
	}
}
