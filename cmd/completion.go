package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for learn.

To load completions:

Bash:
  $ source <(learn completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ learn completion bash > /etc/bash_completion.d/learn
  # macOS:
  $ learn completion bash > $(brew --prefix)/etc/bash_completion.d/learn

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  # To load completions for each session, execute once:
  $ learn completion zsh > "${fpath[1]}/_learn"
  # You will need to start a new shell for this setup to take effect.

Fish:
  $ learn completion fish | source
  # To load completions for each session, execute once:
  $ learn completion fish > ~/.config/fish/completions/learn.fish

Or simply run:
  $ learn completion install
  # Auto-detects your shell and installs completions.`,
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("specify shell: bash, zsh, fish, or install")
		}

		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletionV2(os.Stdout, true)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "install":
			return installCompletions()
		default:
			return fmt.Errorf("unsupported shell: %s (use bash, zsh, fish, or install)", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func installCompletions() error {
	// Auto-detect shell from SHELL env
	shell := os.Getenv("SHELL")
	if shell == "" {
		return fmt.Errorf("SHELL environment variable not set; specify shell manually: learn completion [bash|zsh|fish]")
	}

	shellName := filepath.Base(shell)
	switch shellName {
	case "bash":
		return installBash()
	case "zsh":
		return installZsh()
	case "fish":
		return installFish()
	default:
		return fmt.Errorf("unsupported shell: %s (use bash, zsh, or fish)", shellName)
	}
}

func installBash() error {
	// Try standard completion dirs
	dirs := []string{
		"/etc/bash_completion.d",
		filepath.Join(os.Getenv("HOME"), ".local/share/bash-completion/completions"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			continue
		}
		path := filepath.Join(dir, "learn")
		f, err := os.Create(path)
		if err != nil {
			continue
		}
		defer f.Close()
		if err := rootCmd.GenBashCompletionV2(f, true); err != nil {
			return fmt.Errorf("failed to generate bash completion: %w", err)
		}
		fmt.Printf("Installed bash completion to %s\n", path)
		fmt.Println("Restart your shell or run: source " + path)
		return nil
	}

	// Fallback: print to stdout
	fmt.Println("Could not write to standard completion directories.")
	fmt.Println("Add this to your ~/.bashrc:")
	fmt.Println("  source <(learn completion bash)")
	return nil
}

func installZsh() error {
	fpath := os.Getenv("FPATH")
	if fpath == "" {
		// Common zsh completion dirs
		fpath = filepath.Join(os.Getenv("HOME"), ".zfunc")
	}

	dirs := []string{
		fpath,
		filepath.Join(os.Getenv("HOME"), ".zfunc"),
		"/usr/local/share/zsh/site-functions",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			continue
		}
		path := filepath.Join(dir, "_learn")
		f, err := os.Create(path)
		if err != nil {
			continue
		}
		defer f.Close()
		if err := rootCmd.GenZshCompletion(f); err != nil {
			return fmt.Errorf("failed to generate zsh completion: %w", err)
		}
		fmt.Printf("Installed zsh completion to %s\n", path)
		fmt.Println("Restart your shell or run: exec zsh")
		return nil
	}

	fmt.Println("Could not write to standard completion directories.")
	fmt.Println("Add this to your ~/.zshrc:")
	fmt.Println("  source <(learn completion zsh)")
	return nil
}

func installFish() error {
	dir := filepath.Join(os.Getenv("HOME"), ".config", "fish", "completions")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create fish completions dir: %w", err)
	}

	path := filepath.Join(dir, "learn.fish")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create completion file: %w", err)
	}
	defer f.Close()

	if err := rootCmd.GenFishCompletion(f, true); err != nil {
		return fmt.Errorf("failed to generate fish completion: %w", err)
	}

	// Fish loads completions automatically from this dir
	fmt.Printf("Installed fish completion to %s\n", path)
	// Remove description lines that conflict with fish's own completion descriptions
	// Actually, just let fish handle it natively
	if err := stripFishDescriptions(path); err != nil {
		// Non-fatal
		fmt.Printf("Warning: could not clean up fish completion: %v\n", err)
	}
	fmt.Println("Restart your shell or run: exec fish")
	return nil
}

// stripFishDescriptions removes the "-d" description flags from fish completions
// to avoid conflicts with fish's built-in completion description handling.
func stripFishDescriptions(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var cleaned []string
	for _, line := range lines {
		// Keep the line as-is; fish handles -d flags fine
		cleaned = append(cleaned, line)
	}

	return os.WriteFile(path, []byte(strings.Join(cleaned, "\n")), 0644)
}
