package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"learn/internal/config"
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
		check("rg", hasBinary("rg"))
		check("bat", hasBinary("bat"))
		check("glow", hasBinary("glow"))
		check("EDITOR", getEditor() != "")

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
			// Does the directory actually exist?
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

func hasBinary(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func getEditor() string {
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	if hasBinary("nvim") {
		return "nvim"
	}
	if hasBinary("vim") {
		return "vim"
	}
	if hasBinary("vi") {
		return "vi"
	}
	return ""
}
