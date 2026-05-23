# Go Harness

Go Harness is a minimal terminal-based coding agent written in Go. It uses Google Gemini and exposes a small toolset so the agent can read files, write files, run shell commands, list directories, and optionally search the web through Tavily.

## Requirements

- Go 1.26.2 or newer
- A Gemini API key
- A Tavily API key for web search

## Environment Variables

The program loads variables from a local `.env` file automatically. You can also export the same variables in your shell.

```env
GEMINI_API_KEY=your_gemini_api_key
TAVILY_API_KEY=your_tavily_api_key
```

Required variables:

- `GEMINI_API_KEY`: Gemini API key used to call the Google Generative Language API.
- `TAVILY_API_KEY`: Tavily API key used by the `web_search` tool.

## Run The Program

Install dependencies:

```sh
go mod tidy
```

Create a `.env` file:

```sh
cp .env.example .env
```

If there is no `.env.example`, create `.env` manually with the variables shown above.

Start the agent:

```sh
go run .
```

You should see:

```text
Coding agent listo.
Comandos: :plan on/off, :supervision on/off, :quit
```

Then type a task at the prompt.

## Interactive Commands

- `:plan on`: Ask the agent to propose a plan before running the task.
- `:plan off`: Disable plan confirmation mode.
- `:supervision on`: Ask for confirmation before tools that modify the system run.
- `:supervision off`: Disable supervision prompts.
- `:quit` or `:exit`: Exit the program.

## Available Tools

- `read_file`: Reads a file from disk.
- `write_file`: Writes content to a file.
- `run_command`: Runs a shell command.
- `list_files`: Lists files in a directory.
- `web_search`: Searches the web with Tavily when `TAVILY_API_KEY` is set.

## Notes

- The agent sends requests to `https://generativelanguage.googleapis.com/v1beta/models/{model}:generateContent`.
- Tool calls are printed in the terminal before they run.
- When supervision mode is enabled, `write_file` and `run_command` require confirmation.
