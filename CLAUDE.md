# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CodeGPT is a CLI tool written in Go that generates git commit messages and code reviews using AI. It supports multiple AI providers (OpenAI, Azure OpenAI, Gemini, Anthropic, Ollama, Groq, OpenRouter) and integrates with git prepare-commit-msg hooks.

## Build and Development Commands

```bash
# Build the binary
make build              # outputs to bin/codegpt

# Install globally
make install

# Run all tests with coverage
make test               # or: go test -v -cover -coverprofile coverage.txt ./...

# Run a single test
go test -v -run TestName ./path/to/package

# Format go files (auto-installs golangci-lint if missing)
make fmt

# Lint (auto-installs golangci-lint v2 if missing)
make lint
```

## Development Guidelines

**IMPORTANT: When making changes to this codebase, you MUST:**

1. **Write Tests** - Always write corresponding tests for any new functionality or bug fixes
2. **Pass Linting** - All code must pass `make lint` validation before submission
3. **Format Code** - Always run `make fmt` to ensure consistent coding style

## Architecture

### Core Design Pattern

CodeGPT uses a **provider pattern** where all AI services implement the `core.Generative` interface:

```go
type Generative interface {
    Completion(ctx context.Context, content string) (*Response, error)
    GetSummaryPrefix(ctx context.Context, content string) (*Response, error)
}
```

This allows adding new AI providers without changing core logic.

### Package Structure

- **cmd/** - CLI commands using Cobra

  - `cmd.go` - Root command setup with Viper config initialization (supports env vars with CI platform prefixes like `GITHUB_` and `DRONE_`)
  - `commit.go` - Generate commit messages (`codegpt commit`)
  - `review.go` - Generate code reviews (`codegpt review`)
  - `hook.go` - Git hook management (`codegpt hook install/uninstall`)
  - `provider.go` - Factory functions for instantiating AI providers (`NewOpenAI`, `NewGemini`, `NewAnthropic`)
  - `config.go`, `config_set.go`, `config_list.go` - Configuration management

- **core/** - Core abstractions

  - `openai.go` - `Generative` interface definition and `Response`/`Usage` types
  - `platform.go` - CI/CD platform detection (GitHub Actions, Drone CI)
  - `transport/` - Custom HTTP transport for headers and proxy handling

- **provider/** - AI provider implementations (all implement `core.Generative`)

  - `openai/` - OpenAI-compatible providers (OpenAI, Azure, Groq, OpenRouter, Ollama)
  - `anthropic/` - Anthropic Claude API
  - `gemini/` - Google Gemini (supports both Gemini API and VertexAI backends via `backend` config)

- **git/** - Git operations

  - `git.go` - Diff generation with context lines, file exclusion patterns, commit/amend operations
  - `hook.go` - prepare-commit-msg hook installation/uninstallation
  - `templates/` - Git hook templates

- **prompt/** - Prompt template management

  - `templates/` - Default prompt templates (conventional_commit.tmpl, code_review_file_diff.tmpl, etc.)
  - Supports custom prompt folders via `prompt.folder` config
  - Uses Go templates with custom variables support

- **proxy/** - Network proxy support

  - Handles both SOCKS5 and HTTP/HTTPS proxies
  - Used by provider clients for network requests

- **util/** - Shared utilities
  - `api_key_helper.go` - Dynamic API key retrieval from shell commands with caching (supports password managers, secret services)
  - `api_key_helper_unix.go` / `api_key_helper_windows.go` - Platform-specific subprocess execution
  - `template.go` - Go template rendering with custom functions
  - Cache stored in `$HOME/.config/codegpt/.cache/` with 0600 permissions

### Configuration System

Config file location: `$HOME/.config/codegpt/.codegpt.yaml`
Prompt folder: `$HOME/.config/codegpt/prompt/`
Cache folder: `$HOME/.config/codegpt/.cache/`

**Configuration priority** (for API keys):

1. `{provider}.api_key_helper` or `openai.api_key_helper` (dynamic retrieval)
2. Static config key (`openai.api_key`, `gemini.api_key`)
3. Environment variables (`OPENAI_API_KEY`)

**Key config sections:**

- `openai.*` - Provider settings (provider, api_key, model, base_url, timeout, max_tokens, temperature, etc.)
- `gemini.*` - Gemini-specific settings (api_key, backend, project_id, location)
- `git.*` - Git diff settings (diff_unified, exclude_list)
- `output.lang` - Output language (en, zh-tw, zh-cn, ja)
- `prompt.folder` - Custom prompt template directory

### Adding a New Provider

1. Create package under `provider/<name>/`
2. Define a `Client` struct
3. Implement `core.Generative` interface:
   - `Completion(ctx, content)` - Main completion method
   - `GetSummaryPrefix(ctx, content)` - Generate conventional commit prefix (feat, fix, etc.)
4. Use **options pattern** for configuration (e.g., `WithToken()`, `WithModel()`, `WithTimeout()`)
5. Add factory function in `cmd/provider.go` (e.g., `NewMyProvider(ctx)`) that reads config via Viper
6. Update provider initialization logic in command files to support new provider name
