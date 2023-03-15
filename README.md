# CodeGPT

[![Lint and Testing](https://github.com/appleboy/CodeGPT/actions/workflows/testing.yml/badge.svg?branch=main)](https://github.com/appleboy/CodeGPT/actions/workflows/testing.yml)
[![codecov](https://codecov.io/gh/appleboy/CodeGPT/branch/main/graph/badge.svg)](https://codecov.io/gh/appleboy/CodeGPT)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/CodeGPT)](https://goreportcard.com/report/github.com/appleboy/CodeGPT)

![cover](./images/cover.png)

A CLI written in [Go](https://go.dev) language that writes git commit messages for you using ChatGPT AI (gpt-3.5-turbo model) and automatically installs a [git prepare-commit-msg hook](https://git-scm.com/docs/githooks).

[繁體中文介紹][1]

[1]:https://blog.wu-boy.com/2023/03/writes-git-commit-messages-using-chatgpt/

## Feature

* Support [conventional commits specification](https://www.conventionalcommits.org/en/v1.0.0/).
* Support Git prepare-commit-msg Hook, see the [Git Hooks documentation](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks).
* Support customize generate diffs with n lines of context, default is 3.
* Support exclude file list from git diff command.
* Support translate commit message to another language (support `en`, `zh-tw` or `zh-cn`).
* Support custom network http proxy or socks proxy.

## Installation

Currently, the only supported method of installation on MacOS is [Homebrew](http://brew.sh/). To install `codegpt` via brew:

```sh
brew tap appleboy/tap
brew install codegpt
```

The pre-compiled binaries can be downloaded from [release page](https://github.com/appleboy/CodeGPT/releases).

On linux AMD64

```sh
wget -c https://github.com/appleboy/CodeGPT/releases/download/v0.0.9/CodeGPT-0.0.9-linux-amd64 -O codegpt
```

On macOS (Intel amd64)

```sh
wget -c https://github.com/appleboy/CodeGPT/releases/download/v0.0.9/CodeGPT-0.0.9-darwin-amd64 -O codegpt
```

On macOS (Apple arm64)

```sh
wget -c https://github.com/appleboy/CodeGPT/releases/download/v0.0.9/CodeGPT-0.0.9-darwin-arm64 -O codegpt
```

On Windows (AMD64)

```sh
wget -c https://github.com/appleboy/CodeGPT/releases/download/v0.0.9/CodeGPT-0.0.9-windows-amd64.exe -O codegpt.exe
```

Change the binary permissions to `755` and copy the binary to the system bin directory. Use the `codegpt` command as shown below.

```sh
$ codegpt version
version: v0.0.9 commit: xxxxxxx
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

* **openai.api_key**: generate API key from [openai platform page](https://platform.openai.com/account/api-keys).
* **openai.org_id**: Identifier for this organization sometimes used in API requests. see [organization settings](https://platform.openai.com/account/org-settings).
* **openai.model**: default model is `gpt-3.5-turbo`, you can change to `text-davinci-003` or [other available model list](https://github.com/appleboy/CodeGPT/blob/a75ed831ce30c5c593613b9c0792954586d7f399/openai/openai.go#L16-L29).
* **openai.lang**: default language is `en` and available languages `zh-tw`, `zh-tw`, `ja`.
* **openai.proxy**: http/https client proxy.
* **openai.socks**: socks client proxy.
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

## Reference

* [OpenAI Chat completions documentation](https://platform.openai.com/docs/guides/chat).
* [Introducing ChatGPT and Whisper APIs](https://openai.com/blog/introducing-chatgpt-and-whisper-apis)
