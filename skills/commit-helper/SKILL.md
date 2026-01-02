---
name: generating-commit-messages
description: Generates git commit messages automatically by analyzing your staged changes using AI. Use this skill when you want to save time writing commit messages or need help describing what your code changes do.
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

### Advanced Options

- **Set language**: Use `--lang` to specify output language (en, zh-tw, zh-cn)

  ```bash
  codegpt commit --lang zh-tw --no_confirm
  ```

- **Use specific model**: Override the default model

  ```bash
  codegpt commit --model gpt-4o --no_confirm
  ```

- **Exclude files**: Ignore certain files from the diff analysis

  ```bash
  codegpt commit --exclude_list "*.lock,*.json" --no_confirm
  ```

- **Custom templates**: Format messages according to your team's style

  ```bash
  codegpt commit --template_string "[{{.summarize_prefix}}] {{.summarize_title}}" --no_confirm
  ```

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

### Example 4: Custom template with ticket number

**Input:**

```bash
git add src/api/payment.go
codegpt commit --template_string "{{.summarize_prefix}}(JIRA-123): {{.summarize_title}}" --no_confirm
```

**Output:**

```txt
feat(JIRA-123): integrate payment gateway API
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

**Issue**: Team requires specific commit message format (e.g., with ticket numbers).

**Solution**: Use custom templates:

```bash
codegpt commit --template_string "[{{.summarize_prefix}}](TICKET-123): {{.summarize_title}}"
```

Or save it in config file:

```yaml
git:
  template_string: "[{{.summarize_prefix}}]({{.ticket}}): {{.summarize_title}}"
```

### Multilingual team

**Issue**: Need commit messages in different languages for different repositories.

**Solution**: Set language per command or configure per repository:

```bash
# Per command
codegpt commit --lang zh-cn --no_confirm

# Or set in repository's .codegpt.yaml
output:
  lang: zh-cn
```

### Network proxy required

**Issue**: API calls fail due to corporate firewall.

**Solution**: Configure proxy settings:

```bash
codegpt commit --proxy http://proxy.company.com:8080
# Or for SOCKS proxy
codegpt commit --socks socks5://127.0.0.1:1080
```
