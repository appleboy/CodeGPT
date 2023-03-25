# CodeGPT

[![Lint and Testing](https://github.com/appleboy/CodeGPT/actions/workflows/testing.yml/badge.svg?branch=main)](https://github.com/appleboy/CodeGPT/actions/workflows/testing.yml)
[![codecov](https://codecov.io/gh/appleboy/CodeGPT/branch/main/graph/badge.svg)](https://codecov.io/gh/appleboy/CodeGPT)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/CodeGPT)](https://goreportcard.com/report/github.com/appleboy/CodeGPT)

![cover](./images/cover.png)

A CLI written in [Go](https://go.dev) language that writes git commit messages for you using ChatGPT AI (gpt-3.5-turbo, gpt-4 model) and automatically installs a [git prepare-commit-msg hook](https://git-scm.com/docs/githooks).

[繁體中文介紹][1]

[1]:https://blog.wu-boy.com/2023/03/writes-git-commit-messages-using-chatgpt/

## Feature

* Support [conventional commits specification](https://www.conventionalcommits.org/en/v1.0.0/).
* Support Git prepare-commit-msg Hook, see the [Git Hooks documentation](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks).
* Support customize generate diffs with n lines of context, the default is three.
* Support for excluding files from the git diff command.
* Support commit message translation into another language (support `en`, `zh-tw` or `zh-cn`).
* Support socks proxy or custom network HTTP proxy.
* Support [model lists](https://github.com/appleboy/CodeGPT/blob/bf28f000463cfc6dfa2572df61e1b160c5c680f7/openai/openai.go#L18-L38) like `gpt-4`, `gpt-3.5-turbo` ...etc.
* Support do a brief code review.

![code review](./images/code_review.png)

## Installation

Currently, the only supported method of installation on MacOS is [Homebrew](http://brew.sh/). To install `codegpt` via brew:

```sh
brew tap appleboy/tap
brew install codegpt
```

The pre-compiled binaries can be downloaded from [release page](https://github.com/appleboy/CodeGPT/releases).

On linux AMD64

```sh
wget -c https://github.com/appleboy/CodeGPT/releases/download/v0.1.3/CodeGPT-0.1.3-linux-amd64 -O codegpt
```

On macOS (Intel amd64)

```sh
wget -c https://github.com/appleboy/CodeGPT/releases/download/v0.1.3/CodeGPT-0.1.3-darwin-amd64 -O codegpt
```

On macOS (Apple arm64)

```sh
wget -c https://github.com/appleboy/CodeGPT/releases/download/v0.1.3/CodeGPT-0.1.3-darwin-arm64 -O codegpt
```

On Windows (AMD64)

```sh
wget -c https://github.com/appleboy/CodeGPT/releases/download/v0.1.3/CodeGPT-0.1.3-windows-amd64.exe -O codegpt.exe
```

Change the binary permissions to `755` and copy the binary to the system bin directory. Use the `codegpt` command as shown below.

```sh
$ codegpt version
version: v0.1.3 commit: xxxxxxx
```

## Setup

Please first create your OpenAI API Key. The [OpenAI Platform](https://platform.openai.com/account/api-keys) allows you to generate a new API Key.

![register](./images/register.png)

An environment variable is a variable that is set on your operating system, rather than within your application. It consists of a name and value.We recommend that you set the name of the variable to `OPENAI_API_KEY`.

See the [Best Practices for API Key Safety](https://help.openai.com/en/articles/5112595-best-practices-for-api-key-safety).

```sh
export OPENAI_API_KEY=sk-xxxxxxx
```

or store your API key in custom config file.

```sh
codegpt config set openai.api_key sk-xxxxxxx
```

This will create a `.codegpt.yaml` file in your home directory ($HOME/.config/codegpt/.codegpt.yaml). The following options are available.

* **openai.base_url**: replace the default base URL (`https://api.openai.com/v1`). You can try `https://closeai.deno.dev/v1`. See [justjavac/openai-proxy](https://github.com/justjavac/openai-proxy).
* **openai.api_key**: generate API key from [openai platform page](https://platform.openai.com/account/api-keys).
* **openai.org_id**: Identifier for this organization sometimes used in API requests. see [organization settings](https://platform.openai.com/account/org-settings).
* **openai.model**: default model is `gpt-3.5-turbo`, you can change to `gpt-4` or [other available model list](https://github.com/appleboy/CodeGPT/blob/bf28f000463cfc6dfa2572df61e1b160c5c680f7/openai/openai.go#L18-L38).
* **openai.lang**: default language is `en` and available languages `zh-tw`, `zh-tw`, `ja`.
* **openai.proxy**: http/https client proxy.
* **openai.socks**: socks client proxy.
* **openai.timeout**: default http timeout is `10s` (ten seconds).
* **git.diff_unified**: generate diffs with `<n>` lines of context, default is `3`.
* **git.exclue_list**: exclude file from `git diff` command.

## Usage

There are two methods for generating a commit message using the `codegpt` command. The first is CLI mode, and the second is Git Hook.

### CLI mode

You can call `codegpt` directly to generate a commit message for your staged changes:

```sh
git add <files...>
codegpt commit --preview
```

The commit message is shown below.

```sh
Summarize the commit message use gpt-3.5-turbo model
We are trying to summarize a git diff
We are trying to summarize a title for pull request
================Commit Summary====================

feat: Add preview flag and remove disableCommit flag in commit command and template file.

- Add a `preview` flag to the `commit` command
- Remove the `disbaleCommit` flag from the `prepare-commit-msg` template file

==================================================
Write the commit message to .git/COMMIT_EDITMSG file
```

or translate all git commit messages into a different language (`Traditional Chinese`, `Simplified Chinese` or `Japanese`)

```sh
codegpt commit --lang zh-tw --preview
```

Consider the following outcome:

```sh
Summarize the commit message use gpt-3.5-turbo model
We are trying to summarize a git diff
We are trying to summarize a title for pull request
We are trying to translate a git commit message to Traditional Chinese language
================Commit Summary====================

功能：重構 codegpt commit 命令標記

- 將「codegpt commit」命令新增「預覽」標記
- 從「codegpt commit」命令中移除「--disableCommit」標記

==================================================
Write the commit message to .git/COMMIT_EDITMSG file
```

You can replace the tip of the current branch by creating a new commit. just use `--amend` flag

```sh
codegpt commit --amend
```

## Change commit message template

Default commit message template as following:

```tmpl
{{ .summarize_prefix }}: {{ .summarize_title }}

{{ .summarize_message }}
```

change format with template string using `--template_string` paratemter:

```sh
codegpt commit --preview --template_string \
  "[{{ .summarize_prefix }}]: {{ .summarize_title }}"
```

change format with template file using `--template_file` parameter:

```sh
codegpt commit --preview --template_file your_file_path
```

### Git hook

You can also use the prepare-commit-msg hook to integrate `codegpt` with Git. This allows you to use Git normally and edit the commit message before committing.

#### Install

You want to install the hook in the Git repository:

```sh
codegpt hook install
```

#### Uninstall

You want to remove the hook from the Git repository:

```sh
codegpt hook uninstall
```

Stage your files and commit after installation:

```sh
git add <files...>
git commit
```

`codegpt` will generate the commit message for you and pass it back to Git. Git will open it with the configured editor for you to review/edit it. Then, to commit, save and close the editor!

```sh
$ git commit
Summarize the commit message use gpt-3.5-turbo model
We are trying to summarize a git diff
We are trying to summarize a title for pull request
================Commit Summary====================

Improve user experience and documentation for OpenAI tools

- Add download links for pre-compiled binaries
- Include instructions for setting up OpenAI API key
- Add a CLI mode for generating commit messages
- Provide references for OpenAI Chat completions and ChatGPT/Whisper APIs

==================================================
Write the commit message to .git/COMMIT_EDITMSG file
[main 6a9e879] Improve user experience and documentation for OpenAI tools
 1 file changed, 56 insertions(+)
```

### Code Review

You can use `codegpt` to generate a code review message for your staged changes:

```sh
codegpt review
```

or translate all code review messages into a different language (`Traditional Chinese`, `Simplified Chinese` or `Japanese`)

```sh
codegpt review --lang zh-tw
```

See the following result:

```sh
Code review your changes using gpt-3.5-turbo model
We are trying to review code changes
PromptTokens: 1021, CompletionTokens: 200, TotalTokens: 1221
We are trying to translate core review to Traditional Chinese language
PromptTokens: 287, CompletionTokens: 199, TotalTokens: 486
================Review Summary====================

總體而言，此程式碼修補似乎在增加 Review 指令的功能，允許指定輸出語言並在必要時進行翻譯。以下是需要考慮的潛在問題：

- 輸出語言沒有進行輸入驗證。如果指定了無效的語言代碼，程式可能會崩潰或產生意外結果。
- 此使用的翻譯 API 未指定，因此不清楚是否存在任何安全漏洞。
- 無法處理翻譯 API 調用的錯誤。如果翻譯服

==================================================
```

## Reference

* [OpenAI Chat completions documentation](https://platform.openai.com/docs/guides/chat).
* [Introducing ChatGPT and Whisper APIs](https://openai.com/blog/introducing-chatgpt-and-whisper-apis)
