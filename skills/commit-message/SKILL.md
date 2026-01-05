---
name: commit-message
description: Automatically generates, formats, organizes, and improves git commit messages by analyzing your staged changes using AI. Use this skill when you want to create, write, review, or refine commit messages - whether you need to save time, maintain consistent formatting, or better describe your code changes.
---

# Generating Commit Messages

## Step-by-Step Instructions

### Installation

1. Run the install script to download and set up CodeGPT:

   ```bash
   bash < <(curl -sSL https://raw.githubusercontent.com/appleboy/CodeGPT/main/install.sh)
   ```

2. Configure your AI provider in `~/.config/codegpt/.codegpt.yaml`:

   ```yaml
   openai:
     provider: openai # or: azure, anthropic, gemini, ollama, groq, openrouter
     api_key: your_api_key_here
     model: gpt-4o
   ```

### Basic Usage

1. Stage your changes:

   ```bash
   git add <files>
   ```

2. Generate and commit with AI-generated message:

   ```bash
   codegpt commit --no_confirm
   ```

   **Note**:

   - Commit messages are generated in **English by default**. Use `--lang` to specify a different language.
   - You don't need to manually run `git diff` to review your changes. CodeGPT automatically reads the staged changes and analyzes them to generate an appropriate commit message.

### Advanced Options

- **Set language**: Use `--lang` to specify output language (default: en)

  ```bash
  codegpt commit --lang zh-tw --no_confirm  # For Traditional Chinese
  codegpt commit --lang zh-cn --no_confirm  # For Simplified Chinese
  ```

  To change the default language permanently:

  ```bash
  codegpt config set output.lang en  # or zh-tw, zh-cn
  ```

- **Use specific model**: Override the default model

  ```bash
  codegpt commit --model gpt-4o --no_confirm
  ```

- **Exclude files**: Ignore certain files from the diff analysis

  ```bash
  codegpt commit --exclude_list "*.lock,*.json" --no_confirm
  ```

- **Custom templates**: Format messages according to your team's style, including Jira issue tracking

  ```bash
  # Basic custom format
  codegpt commit --template_string "[{{.summarize_prefix}}] {{.summarize_title}}" --no_confirm

  # With Jira issue number using template variables
  codegpt commit --template_vars "JIRA_NUM=GAIA-2704" \
    --template_string "{{.summarize_prefix}}{{if .JIRA_NUM}}({{.JIRA_NUM}}){{end}}: {{.summarize_title}}\n\n{{.summarize_message}}" \
    --no_confirm
  ```

  Available built-in template variables:

  - `{{.summarize_prefix}}` - Conventional commit type (feat, fix, docs, etc.)
  - `{{.summarize_title}}` - The commit title
  - `{{.summarize_message}}` - The full commit message body

  You can define custom variables with `--template_vars "KEY=VALUE"` and use them as `{{.KEY}}`

- **Amend commit**: Update the previous commit message

  ```bash
  codegpt commit --amend --no_confirm
  ```

## Examples of Inputs and Outputs

### Example 1: Adding a new feature

**Input:**

```bash
# After making changes to add user authentication
git add src/auth.go src/middleware.go
codegpt commit --no_confirm
```

**Output:**

```txt
feat: add user authentication middleware

Implement JWT-based authentication system with login and token validation
middleware for protecting API endpoints.
```

### Example 2: Fixing a bug

**Input:**

```bash
# After fixing a null pointer error
git add src/handlers/user.go
codegpt commit --no_confirm
```

**Output:**

```txt
fix: resolve null pointer exception in user handler

Add nil checks before accessing user object properties to prevent crashes
when processing requests with missing user data.

[Preview shown, waiting for confirmation...]
```

### Example 3: Chinese language commit

**Input:**

```bash
git add docs/README.md
codegpt commit --lang zh-tw --no_confirm
```

**Output:**

```txt
docs: 更新專案說明文件

新增安裝步驟說明以及使用範例，讓新使用者能夠快速上手。
```

### Example 4: Jira issue tracking integration

**Input:**

```bash
# Stage your changes
git add Dockerfile docker-compose.yml

# Generate commit with Jira issue number using template variables
codegpt commit --template_vars "JIRA_NUM=GAIA-2704" \
  --template_string "{{.summarize_prefix}}{{if .JIRA_NUM}}({{.JIRA_NUM}}){{end}}: {{.summarize_title}}

{{.summarize_message}}" \
  --no_confirm
```

**Output:**

```txt
feat(GAIA-2704): update trivy scan for cicd

- Add appuser (uid=1000) to run containers as non-root user
- Fix trivy scan DS002: Image user should not be 'root'
- Create /.app directory with proper ownership for externalimage
```

**Tip**: Save the template in config file to avoid typing it every time:

```bash
# Set the template once
codegpt config set git.template_string "{{.summarize_prefix}}{{if .JIRA_NUM}}({{.JIRA_NUM}}){{end}}: {{.summarize_title}}

{{.summarize_message}}"

# Then just provide the Jira number when committing
codegpt commit --template_vars "JIRA_NUM=GAIA-2704" --no_confirm
```

### Example 5: Excluding lock files

**Input:**

```bash
git add .
codegpt commit --exclude_list "package-lock.json,yarn.lock,go.sum" --no_confirm
```

**Output:**

```txt
refactor: reorganize project structure

Move utility functions into separate packages and update import paths
throughout the codebase for better modularity.

(Lock files excluded from analysis)
```

## Common Edge Cases

### No staged changes

**Issue**: Running `codegpt commit` without staging any changes.

**Solution**: Stage your changes first:

```bash
git add <files>
codegpt commit
```

### API timeout for large diffs

**Issue**: Large changesets may cause API timeouts.

**Solution**: Increase timeout or commit changes in smaller batches:

```bash
codegpt commit --timeout 60s --no_confirm
```

### Generated files in diff

**Issue**: Lock files or generated code affecting commit message quality.

**Solution**: Exclude these files from analysis:

```bash
codegpt commit --exclude_list "package-lock.json,yarn.lock,*.min.js,dist/*" --no_confirm
```

### API key not configured

**Issue**: Error message about missing API key.

**Solution**: Set up your API key in the config file:

```bash
codegpt config set openai.api_key "your-api-key-here"
```

### Custom commit format required

**Issue**: Team requires specific commit message format with Jira issue tracking numbers.

**Solution**: Use `--template_vars` with custom templates.

**Basic usage:**

```bash
codegpt commit --template_vars "JIRA_NUM=GAIA-2704" \
  --template_string "{{.summarize_prefix}}{{if .JIRA_NUM}}({{.JIRA_NUM}}){{end}}: {{.summarize_title}}

{{.summarize_message}}" \
  --no_confirm

# Result: feat(GAIA-2704): update trivy scan for cicd
```

**Save template in config for convenience:**

```bash
# Set the template once
codegpt config set git.template_string "{{.summarize_prefix}}{{if .JIRA_NUM}}({{.JIRA_NUM}}){{end}}: {{.summarize_title}}

{{.summarize_message}}"

# Then just provide the Jira number
codegpt commit --template_vars "JIRA_NUM=GAIA-2704" --no_confirm
```

Or edit `~/.config/codegpt/.codegpt.yaml`:

```yaml
git:
  template_string: "{{.summarize_prefix}}{{if .JIRA_NUM}}({{.JIRA_NUM}}){{end}}: {{.summarize_title}}\n\n{{.summarize_message}}"
```

**Create a shell function for convenience:**

Add to your `~/.bashrc` or `~/.zshrc`:

```bash
function commit() {
  if [ -z "$1" ]; then
    echo "Usage: commit <JIRA-NUM>"
    return 1
  fi
  codegpt commit --template_vars "JIRA_NUM=$1" --no_confirm
}

# Usage: commit GAIA-2704
```

**Auto-extract from git branch name (advanced):**

If your branch follows naming convention like `feature/GAIA-2704-add-security-scan`:

```bash
# Add to ~/.bashrc or ~/.zshrc
function commit-auto() {
  local branch=$(git rev-parse --abbrev-ref HEAD)
  local jira_num=$(echo "$branch" | grep -oE '[A-Z]+-[0-9]+' | head -1)

  if [ -z "$jira_num" ]; then
    echo "No Jira issue found in branch: $branch"
    codegpt commit --no_confirm
  else
    echo "Detected Jira issue: $jira_num"
    codegpt commit --template_vars "JIRA_NUM=$jira_num" --no_confirm
  fi
}

# Usage: commit-auto (auto-detects GAIA-2704 from branch name)
```

### Network proxy required

**Issue**: API calls fail due to corporate firewall.

**Solution**: Configure proxy settings:

```bash
codegpt commit --proxy http://proxy.company.com:8080
# Or for SOCKS proxy
codegpt commit --socks socks5://127.0.0.1:1080
```
