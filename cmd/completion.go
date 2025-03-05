package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var CompletionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `Generate shell completion scripts for CodeGPT CLI.
	
To load completions:

Bash:
	$ source <(codegpt completion bash)
	# Or save it to a file and source it:
	$ codegpt completion bash > ~/.codegpt-completion.bash
	$ echo 'source ~/.codegpt-completion.bash' >> ~/.bashrc

Zsh:
	$ source <(codegpt completion zsh)
	# Or save it to a file in your $fpath:
	$ codegpt completion zsh > "${fpath[1]}/_codegpt"

Fish:
	$ codegpt completion fish > ~/.config/fish/completions/codegpt.fish

PowerShell:
	PS> codegpt completion powershell > codegpt.ps1
	PS> . ./codegpt.ps1`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			_ = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			_ = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			_ = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}
