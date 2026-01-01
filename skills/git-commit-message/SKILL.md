---
name: git commit message
description: This skill provides AI-powered git commit message generation using the `codegpt commit` command. It analyzes git diffs and automatically generates conventional commit messages in multiple languages.
---

# Git Commit Message

## Description

The git-commit skill leverages various AI providers (OpenAI, Anthropic, Gemini, Ollama, Groq, OpenRouter) to automatically generate meaningful commit messages that follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.

## Installation

Run the install script to automatically download and set up the latest release:

```sh
bash < <(curl -sSL https://raw.githubusercontent.com/appleboy/CodeGPT/main/install.sh)
```

## Usage

### Basic Usage

Generate a commit message for staged changes:

```bash
codegpt commit
```

### Common Options

```bash
# Preview commit message before committing
codegpt commit --preview

# Skip confirmation prompts
codegpt commit --no_confirm

# Set output language (en, zh-tw, zh-cn)
codegpt commit --lang zh-tw

# Use specific AI model
codegpt commit --model gpt-4o

# Amend previous commit
codegpt commit --amend

# Display prompt only (no API call)
codegpt commit --prompt_only

# Write to specific output file
codegpt commit --file /path/to/commit-msg

# Customize diff context lines
codegpt commit --diff_unified 5

# Exclude specific files from diff
codegpt commit --exclude_list "*.lock,*.json"

# Use custom template file
codegpt commit --template_file ./my-template.tmpl

# Use inline template string
codegpt commit --template_string "{{.summarize_prefix}}: {{.summarize_title}}"

# Set API timeout
codegpt commit --timeout 60s

# Configure network proxy
codegpt commit --proxy http://proxy.example.com:8080
codegpt commit --socks socks5://127.0.0.1:1080
```

### Template Variables

When using custom templates, the following variables are available:

- `{{.summarize_prefix}}` - Conventional commit prefix (feat, fix, docs, etc.)
- `{{.summarize_title}}` - Brief commit title
- `{{.summarize_message}}` - Detailed commit message body

You can also provide custom variables:

```bash
codegpt commit --template_vars "author=John,ticket=PROJ-123"
codegpt commit --template_vars_file ./vars.env
```

## Workflow

1. **Diff Analysis**: Analyzes staged changes using `git diff`
2. **Summarization**: AI generates a summary of the changes
3. **Title Generation**: Creates a concise commit title
4. **Prefix Detection**: Determines appropriate conventional commit prefix
5. **Message Composition**: Combines all elements into a formatted message
6. **Translation** (optional): Translates to target language if specified
7. **Preview & Confirmation**: Shows message for review and optional editing
8. **Commit**: Records changes to the repository

## Configuration

Configure via `~/.config/codegpt/.codegpt.yaml`:

```yaml
openai:
  provider: openai  # or: azure, anthropic, gemini, ollama, groq, openrouter
  api_key: your_api_key_here
  model: gpt-4o
  timeout: 30s

git:
  diff_unified: 3
  exclude_list: []
  template_file: ""
  template_string: ""

output:
  lang: en  # or: zh-tw, zh-cn
  file: ""
```

## Examples

### Example 1: Basic commit with preview

```bash
# Stage your changes
git add .

# Generate and preview commit message
codegpt commit --preview
```

### Example 2: Chinese commit message

```bash
codegpt commit --lang zh-tw --model gpt-4o
```

### Example 3: Custom template

```bash
codegpt commit \
  --template_string "[{{.summarize_prefix}}] {{.summarize_title}}" \
  --template_vars "ticket=PROJ-123"
```

### Example 4: With file exclusions

```bash
codegpt commit \
  --exclude_list "package-lock.json,yarn.lock,go.sum" \
  --preview
```

## Tips

1. **Stage Changes First**: Always run `git add` before using `codegpt commit`
2. **Use Preview Mode**: Review messages with `--preview` before committing
3. **Customize Templates**: Create templates that match your team's commit style
4. **Set Default Language**: Configure your preferred language in config file
5. **Exclude Generated Files**: Use `--exclude_list` to ignore lock files and generated code
