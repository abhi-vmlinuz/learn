package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"learn/internal/config"
	"learn/internal/git"

	"github.com/spf13/cobra"
)

var logLimit int

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show git history of notes",
	Long:  "Display recent git commits affecting your notes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if !git.IsRepo(cfg.Repo.Root) {
			return fmt.Errorf("not a git repository: %s", cfg.Repo.Root)
		}

		limit := fmt.Sprintf("-%d", logLimit)
		gitCmd := exec.Command("git", "log", limit, "--oneline", "--stat", "--no-color")
		gitCmd.Dir = cfg.Repo.Root
		out, err := gitCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git log failed: %s", strings.TrimSpace(string(out)))
		}

		fmt.Print(string(out))
		return nil
	},
}

func init() {
	logCmd.Flags().IntVar(&logLimit, "limit", 10, "Number of commits to show")
	rootCmd.AddCommand(logCmd)
}
