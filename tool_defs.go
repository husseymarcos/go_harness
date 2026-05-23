package main

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
