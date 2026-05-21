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
		tool("read_file", "Lee el contenido de un archivo.", map[string]any{
			"path": stringParam("Path del archivo a leer."),
		}, []string{"path"}),
		tool("write_file", "Escribe contenido en un archivo, reemplazando el contenido actual.", map[string]any{
			"path":    stringParam("Path del archivo a escribir."),
			"content": stringParam("Contenido completo del archivo."),
		}, []string{"path", "content"}),
		tool("run_command", "Ejecuta un comando de terminal y devuelve stdout y stderr.", map[string]any{
			"command": stringParam("Comando completo a ejecutar."),
		}, []string{"command"}),
		tool("list_files", "Lista los archivos en un directorio.", map[string]any{
			"path": stringParam("Directorio a listar. Usar . si no se especifica."),
		}, []string{"path"}),
		tool("web_search", "Busca información en la web con Tavily.", map[string]any{
			"query": stringParam("Consulta de búsqueda."),
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
				"type":                 "object",
				"properties":           properties,
				"required":             required,
				"additionalProperties": false,
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
