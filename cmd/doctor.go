package cmd

import (
	"fmt"
	"os"

	"learn/internal/config"
	"learn/internal/editor"
	"learn/internal/fzf"
	"learn/internal/git"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check environment health",
	Long:  "Verify all dependencies and configuration are properly set up.",
	RunE: func(cmd *cobra.Command, args []string) error {
		allGood := true

		check := func(name string, ok bool) {
			if ok {
				fmt.Printf("  ✓ %s\n", name)
			} else {
				fmt.Printf("  ✗ %s\n", name)
				allGood = false
			}
		}

		// External tools
		check("git", git.IsAvailable())
		check("fzf", fzf.IsAvailable())
		check("rg", editor.HasBinary("rg"))
		check("bat", editor.HasBinary("bat"))
		check("glow", editor.HasBinary("glow"))
		check("tdf", editor.HasBinary("tdf"))
		check("wkhtmltopdf", editor.HasBinary("wkhtmltopdf"))
		check("EDITOR", editor.GetEditor() != "")

		// Config file
		cfgPath := config.ConfigPath()
		_, err := os.Stat(cfgPath)
		check("config file", err == nil)

		// Repository: config loads, root exists, root is a git repo
		cfg, err := config.Load()
		if err != nil {
			check("repository", false)
			fmt.Printf("         run 'learn init' in your notes directory\n")
		} else {
			info, err := os.Stat(cfg.Repo.Root)
			if err != nil || !info.IsDir() {
				check("repository", false)
				fmt.Printf("         root not found: %s\n", cfg.Repo.Root)
				fmt.Printf("         run 'learn init' to reinitialize\n")
			} else if !git.IsRepo(cfg.Repo.Root) {
				check("repository", false)
				fmt.Printf("         root exists but is not a git repo: %s\n", cfg.Repo.Root)
				fmt.Printf("         run 'git init' in that directory\n")
			} else {
				check("repository", true)
				fmt.Printf("         %s\n", cfg.Repo.Root)
			}
		}

		fmt.Println()
		if allGood {
			fmt.Println("Repository healthy")
		} else {
			fmt.Println("Some checks failed. Fix the issues above.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
