# CodeGPT

[![Lint and Testing](https://github.com/appleboy/CodeGPT/actions/workflows/testing.yml/badge.svg?branch=main)](https://github.com/appleboy/CodeGPT/actions/workflows/testing.yml)
[![codecov](https://codecov.io/gh/appleboy/CodeGPT/branch/main/graph/badge.svg)](https://codecov.io/gh/appleboy/CodeGPT)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/CodeGPT)](https://goreportcard.com/report/github.com/appleboy/CodeGPT)

![cover](./images/cover.png)

A CLI written in [Golang](https://go.dev) that writes your git commit messages for you with ChatGPT AI (`gpt-3.5-turbo` model) and install a [git prepare-commit-msg hook](https://git-scm.com/docs/githooks) automatically.

## Installation

The pre-compiled binaries can be downloaded from [release page](https://github.com/appleboy/CodeGPT/releases).

On linux AMD64

```sh
wget -c https://github.com/appleboy/CodeGPT/releases/download/v0.0.1/CodeGPT-0.0.2-linux-amd64 -O codegpt
```

On macOS (Intel amd64)

```sh
wget -c https://github.com/appleboy/CodeGPT/releases/download/v0.0.1/CodeGPT-0.0.2-darwin-amd64 -O codegpt
```

On macOS (Apple arm64)

```sh
wget -c https://github.com/appleboy/CodeGPT/releases/download/v0.0.1/CodeGPT-0.0.2-darwin-arm64 -O codegpt
```

On Windows (AMD64)

```sh
wget -c https://github.com/appleboy/CodeGPT/releases/download/v0.0.2/CodeGPT-0.0.2-windows-amd64.exe -O codegpt.exe
```

Change the binary permission to `755` and move to system bin directory. Try the `codegpt` command as below

```sh
$ codegpt version
version: v0.0.2 commit: 500ae35
```

## Setup

Please set up your OpenAI API Key first. You can create a new API Key from the [OpenAI Platform](https://platform.openai.com/account/api-keys).

> Note: If you haven't already, you'll have to create an account and set up billing.

```sh
codegpt config set openai.api_key sk-xxxxxxx
```

This will create a `.codegpt.yaml` file in your home directory ($HOME/.config/codegpt/.codegpt.yaml).

## Usage

### CLI mode

You can call `codegpt` directly to generate a commit message for your staged changes:

```sh
git add <files...>
codegpt commit
```

## Reference

* [OpenAI Chat completions documentation](https://platform.openai.com/docs/guides/chat).
* [Introducing ChatGPT and Whisper APIs](https://openai.com/blog/introducing-chatgpt-and-whisper-apis)
