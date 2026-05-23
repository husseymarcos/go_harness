# Agent Instructions for Go Harness

## Project Basics
- Minimal Go coding agent using Google Gemini via the `Provider` interface.
- Go 1.26.2, single module. Run with `go run .`, build with `go build .`.
- Only dependency: `github.com/joho/godotenv`.

## Setup
- Create a `.env` file (gitignored) with `GEMINI_API_KEY=...`.
- `TAVILY_API_KEY` is **optional** — `web_search` degrades gracefully at runtime if missing.

## Architecture (Important)
- `llm.go`: Generic `Message`, `ToolCall`, `ToolSchema` types and the `Provider` interface. **Keep this file tiny.** It intentionally has no JSON tags; serialization lives in provider-specific code.
- `gemini.go`: Implements `Provider`. Contains **all** Gemini-specific DTOs and mapping logic. If you ever switch LLM providers, add a new file implementing `Provider` and only change `main.go`.
- `agent.go`: Core agent loop. `Agent` holds a `Provider`, not raw API credentials.
- `tools.go`: Tool implementations (`read_file`, `write_file`, `run_command`, `list_files`, `web_search`).
- `tool_defs.go`: Schema definitions returned as `[]ToolSchema`.

## Code Conventions
- Do not add JSON tags to types in `llm.go`. Wire-format serialization is the provider's job.
- Keep `llm.go` generic. No provider-specific types or logic there.
- Tool schemas are flat `ToolSchema` structs; the old `ToolDef` wrapper was deleted.

## Running / Testing
- There are no tests yet. Verify with `go build .` and `go vet ./...`.
- Entry point: `main.go`. Model (`gemini-2.0-flash`) is hardcoded there.

## Git
- Do **not** commit `go_harness` binary or `.env`.
