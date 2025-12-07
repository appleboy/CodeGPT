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

## Architecture

### Package Structure

- **cmd/** - CLI commands using Cobra
  - `cmd.go` - Root command setup and config initialization
  - `commit.go` - Generate commit messages (`codegpt commit`)
  - `review.go` - Generate code reviews (`codegpt review`)
  - `hook.go` - Git hook management (`codegpt hook install/uninstall`)
  - `config.go`, `config_set.go`, `config_list.go` - Configuration management

- **core/** - Core abstractions
  - `openai.go` - `Generative` interface that all providers implement
  - `platform.go` - CI/CD platform detection

- **provider/** - AI provider implementations (all implement `core.Generative`)
  - `openai/` - OpenAI, Azure, Groq, OpenRouter, Ollama (OpenAI-compatible)
  - `anthropic/` - Anthropic Claude
  - `gemini/` - Google Gemini (supports both Gemini API and VertexAI backends)

- **git/** - Git operations
  - `git.go` - Diff generation, commit operations
  - `hook.go` - prepare-commit-msg hook installation

- **prompt/** - Prompt management
  - Templates stored in `prompt/templates/`
  - Supports custom prompt folders via config

- **proxy/** - Network proxy support (SOCKS and HTTP)

- **util/** - Utilities including Go template rendering

### Configuration

Config file: `$HOME/.config/codegpt/.codegpt.yaml`
Prompt folder: `$HOME/.config/codegpt/prompt/`

Key config keys: `openai.provider`, `openai.api_key`, `openai.model`, `openai.base_url`, `git.diff_unified`, `output.lang`

### Adding a New Provider

1. Create package under `provider/`
2. Implement `core.Generative` interface (`Completion` and `GetSummaryPrefix` methods)
3. Add options pattern for configuration
4. Wire up in `cmd/provider.go`
