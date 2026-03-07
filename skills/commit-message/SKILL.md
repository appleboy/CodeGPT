---
name: commit-message
description: Generate a conventional commit message by analyzing staged git changes. Use when the user wants to create, write, or generate a git commit message from their current staged diff.
---

# Generate Commit Message

Generate a conventional commit message from staged git changes following a structured prompt pipeline.

## Steps

### 1. Stage changes and get the diff

If there are modified files from the current session that haven't been staged yet, run `git add` on those files first to include them in the staged changes.

Then get the staged diff:

```bash
git diff --staged
```

If the diff is empty after this, inform the user that there are no staged changes and stop.

### 2. Analyze the diff

Produce a bullet-point summary of the changes. Follow these rules:

- A line starting with `+` means it was added, `-` means deleted. Lines with neither are context.
- Write every summary comment as a bullet point starting with `-`.
- Do not include file names as part of the comment.
- Do not use `[` or `]` characters in the summary.
- Do not include comments copied from the code.
- Write only the most important comments. When in doubt, write fewer comments.
- Readability is top priority.

Example summary comments for reference (do not copy verbatim):

```
- Increase the number of returned recordings from 10 to 100
- Correct a typo in the GitHub Action name
- Relocate the octokit initialization to a separate file
- Implement an OpenAI API endpoint for completions
```

### 3. Generate the commit title

From the summary, write a single-line commit title:

- Use imperative tense following the kernel git commit style guide.
- Write a high-level title that captures a single specific theme.
- Do not repeat the file summaries or list individual changes.
- No more than 60 characters.
- Lowercase the first character.
- Remove any trailing period.

### 4. Determine the prefix and scope

**Prefix** — choose exactly one label based on the summary:

- `build`: Changes that affect the build system or external dependencies
- `chore`: Updating libraries, copyrights, or other repo settings, includes updating dependencies
- `ci`: Changes to CI configuration files and scripts
- `docs`: Non-code changes, such as fixing typos or adding new documentation
- `feat`: Introduces a new feature to the codebase
- `fix`: Patches a bug in the codebase
- `perf`: A code change that improves performance
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `style`: Changes that do not affect the meaning of the code (white-space, formatting, etc.)
- `test`: Adding missing tests or correcting existing tests

**Scope** — identify the module or package scope from the changed files:

- Look at the file paths in the diff to determine which module, package, or component is affected.
- If all changes are within a single module/package/directory, use that as the scope (e.g., `model`, `git`, `prompt`, `cmd`, `provider`).
- Use the most specific common directory or package name. For example, changes only in `provider/openai/` should use `openai`, not `provider`.
- If changes span multiple modules, pick the one most central to the change's purpose.
- Scope is **required** — always include one.
- Keep the scope short — a single lowercase word.

### 5. Create the commit

Format the commit message as:

```
<prefix>(<scope>): <title>

<summary>
```

Show the formatted message to the user and ask for confirmation before running `git commit`.
