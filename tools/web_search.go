package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type webSearchArgs struct {
	Query string `json:"query"`
}

func webSearch(rawArgs string) string {
	var args webSearchArgs
	if err := parseArgs(rawArgs, &args); err != nil {
		return err.Error()
	}

	apiKey := os.Getenv("TAVILY_API_KEY")
	if apiKey == "" {
		return "Missing TAVILY_API_KEY. web_search is unavailable."
	}

	body, _ := json.Marshal(map[string]any{
		"api_key":      apiKey,
		"query":        args.Query,
		"max_results":  5,
		"include_raw":  false,
		"search_depth": "basic",
	})

	res, err := http.Post("https://api.tavily.com/search", "application/json", bytes.NewReader(body))
	if err != nil {
		return "Error searching the web: " + err.Error()
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return "Error reading Tavily response: " + err.Error()
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Sprintf("Tavily returned status %d: %s", res.StatusCode, string(raw))
	}
	return string(raw)
}
