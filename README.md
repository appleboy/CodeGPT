# CodeGPT

[![Lint and Testing](https://github.com/appleboy/CodeGPT/actions/workflows/testing.yml/badge.svg?branch=main)](https://github.com/appleboy/CodeGPT/actions/workflows/testing.yml)
[![codecov](https://codecov.io/gh/appleboy/CodeGPT/branch/main/graph/badge.svg)](https://codecov.io/gh/appleboy/CodeGPT)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/CodeGPT)](https://goreportcard.com/report/github.com/appleboy/CodeGPT)

![cover](./images/cover.png)

A CLI written in [Go](https://go.dev) language that writes git commit messages or do a code review brief for you using ChatGPT AI (gpt-3.5-turbo, gpt-4 model) and automatically installs a [git prepare-commit-msg hook](https://git-scm.com/docs/githooks).

* [繁體中文介紹][1]
* [繁體中文影片][2]

[1]:https://blog.wu-boy.com/2023/03/writes-git-commit-messages-using-chatgpt/
[2]:https://www.youtube.com/watch?v=4Yei_t6eMZU

![flow](./images/flow.svg)

## Feature

* Support [Azure OpenAI Service](https://azure.microsoft.com/en-us/products/cognitive-services/openai-service) or [OpenAI API](https://platform.openai.com/docs/api-reference).
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

Install from [Homebrew](http://brew.sh/) on MacOS

```sh
brew tap appleboy/tap
brew install codegpt
```

Install from [Chocolatey](https://chocolatey.org/install) on Windows

```sh
choco install codegpt
```

The pre-compiled binaries can be downloaded from [release page](https://github.com/appleboy/CodeGPT/releases).Change the binary permissions to `755` and copy the binary to the system bin directory. Use the `codegpt` command as shown below.

```sh
$ codegpt version
version: v0.4.3 commit: xxxxxxx
```

Install from source code:

```sh
go install github.com/appleboy/CodeGPT/cmd/codegpt@latest
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

* **openai.base_url**: replace the default base URL (`https://api.openai.com/v1`).
* **openai.api_key**: generate API key from [openai platform page](https://platform.openai.com/account/api-keys).
* **openai.org_id**: Identifier for this organization sometimes used in API requests. see [organization settings](https://platform.openai.com/account/org-settings). only for `openai` service.
* **openai.model**: default model is `gpt-3.5-turbo`, you can change to `gpt-4` or [other available model list](https://github.com/appleboy/CodeGPT/blob/bf28f000463cfc6dfa2572df61e1b160c5c680f7/openai/openai.go#L18-L38).
* **openai.proxy**: http/https client proxy.
* **openai.socks**: socks client proxy.
* **openai.timeout**: default http timeout is `10s` (ten seconds).
* **openai.max_tokens**: default max tokens is `300`. see reference [max_tokens](https://platform.openai.com/docs/api-reference/completions/create#completions/create-max_tokens).
* **openai.temperature**: default temperature is `1`. see reference [temperature](https://platform.openai.com/docs/api-reference/completions/create#completions/create-temperature).
* **git.diff_unified**: generate diffs with `<n>` lines of context, default is `3`.
* **git.exclude_list**: exclude file from `git diff` command.
* **openai.provider**: default service provider is `openai`, you can change to `azure`.
* **openai.model_name**: model deployment name (for azure).
* **output.lang**: default language is `en` and available languages `zh-tw`, `zh-cn`, `ja`.
* **openai.top_p**: default top_p is `1.0`. see reference [top_p](https://platform.openai.com/docs/api-reference/completions/create#completions/create-top_p).
* **openai.frequency_penalty**: default frequency_penalty is `0.0`. see reference [frequency_penalty](https://platform.openai.com/docs/api-reference/completions/create#completions/create-frequency_penalty).
* **openai.presence_penalty**: default presence_penalty is `0.0`. see reference [presence_penalty](https://platform.openai.com/docs/api-reference/completions/create#completions/create-presence_penalty).

### How to change to Azure OpenAI Service

Please get the `API key`, `Endpoint` and `Model deployments` list from Azure Resource Management Portal on left menu.

![azure01](./images/azure_01.png)

![azure02](./images/azure_02.png)

Update your config file.

```sh
codegpt config set openai.provider azure
codegpt config set openai.base_url https://xxxxxxxxx.openai.azure.com/
codegpt config set openai.api_key xxxxxxxxxxxxxxxx
codegpt config set openai.model_name xxxxx-gpt-35-turbo
```

### How to change to Groq API Service

Please get the `API key` from Groq API Service, please vist [here][31]. Update the `base_url` and `api_key` in your config file.

```sh
codegpt config set openai.provider openai
codegpt config set openai.base_url https://api.groq.com/openai/v1
codegpt config set openai.api_key gsk_xxxxxxxxxxxxxx
```

Support the [following models][32]:

1. LLaMA2-70b (Meta) **recommended**
2. Mixtral-8x7b (Mistral)
3. Gemma-7b-it (Google)

[31]: https://console.groq.com/keys
[32]: https://console.groq.com/docs/models

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

change format with template string using `--template_string` parameter:

```sh
codegpt commit --preview --template_string \
  "[{{ .summarize_prefix }}]: {{ .summarize_title }}"
```

change format with template file using `--template_file` parameter:

```sh
codegpt commit --preview --template_file your_file_path
```

Add custom variable to git commit message template:

```sh
{{ .summarize_prefix }}: {{ .summarize_title }}

{{ .summarize_message }}

{{ if .JIRA_URL }}{{ .JIRA_URL }}{{ end }}
```

Add custom variable to git commit message template using `--template_vars` parameter:

```sh
codegpt commit --preview --template_file your_file_path --template_vars \
  JIRA_URL=https://jira.example.com/ABC-123
```

Load custom variable from file using `--template_vars_file` parameter:

```sh
codegpt commit --preview --template_file your_file_path --template_vars_file your_file_path
```

See the `template_vars_file` format as following:

```env
JIRA_URL=https://jira.example.com/ABC-123
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

another php example code:

```php
<?php
if( isset( $_POST[ 'Submit' ]  ) ) {
  // Get input
  $target = $_REQUEST[ 'ip' ];
  // Determine OS and execute the ping command.
  if( stristr( php_uname( 's' ), 'Windows NT' ) ) {
    // Windows
    $cmd = shell_exec( 'ping  ' . $target );
  }
  else {
    // *nix
    $cmd = shell_exec( 'ping  -c 4 ' . $target );
  }
  // Feedback for the end user
  $html .= "<pre>{$cmd}</pre>";
}
?>
```

code review result:

```sh
================Review Summary====================

Code review:

1. Security: The code is vulnerable to command injection attacks as the user input is directly used in the shell_exec() function. An attacker can potentially execute malicious commands on the server by injecting them into the 'ip' parameter.
2. Error handling: There is no error handling in the code. If the ping command fails, the error message is not displayed to the user.
3. Input validation: There is no input validation for the 'ip' parameter. It should be validated to ensure that it is a valid IP address or domain name.
4. Cross-platform issues: The code assumes that the server is either running Windows or *nix operating systems. It may not work correctly on other platforms.

Suggestions for improvement:

1. Use escapeshellarg() function to sanitize the user input before passing it to shell_exec() function to prevent command injection.
2. Implement error handling to display error messages to the user if the ping command fails.
3. Use a regular expression to validate the 'ip' parameter to ensure that it is a valid IP address or domain name.
4. Use a more robust method to determine the operating system, such as the PHP_OS constant, which can detect a wider range of operating systems.

==================================================
```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=appleboy/codegpt&type=Date)](https://star-history.com/#appleboy/codegpt&Date)

## Reference

* [OpenAI Chat completions documentation](https://platform.openai.com/docs/guides/chat).
* [Introducing ChatGPT and Whisper APIs](https://openai.com/blog/introducing-chatgpt-and-whisper-apis)
