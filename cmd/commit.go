package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"learn/internal/config"
	"learn/internal/git"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit [message]",
	Short: "Git add and commit notes",
	Long:  "Stage all modified notes and commit with a message.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		repoRoot := cfg.Repo.Root

		// Get modified/untracked files
		files, err := git.Status(repoRoot)
		if err != nil {
			return err
		}

		// Filter to markdown files only
		var mdFiles []string
		for _, f := range files {
			if strings.HasSuffix(f, ".md") {
				mdFiles = append(mdFiles, f)
			}
		}

		if len(mdFiles) == 0 {
			return fmt.Errorf("no modified or untracked notes found")
		}

		// Display summary
		fmt.Println("Modified:")
		for _, f := range mdFiles {
			rel, _ := filepath.Rel(repoRoot, f)
			if rel == "" {
				rel = f
			}
			fmt.Printf("  %s\n", rel)
		}
		fmt.Println()

		// Get commit message
		var message string
		if len(args) > 0 {
			message = args[0]
		} else {
			fmt.Print("Commit message: ")
			fmt.Scanln(&message)
			if message == "" {
				return fmt.Errorf("commit message cannot be empty")
			}
		}

		// Stage all and commit
		if err := git.AddAll(repoRoot); err != nil {
			return fmt.Errorf("git add failed: %w", err)
		}

		if err := git.Commit(repoRoot, message); err != nil {
			return fmt.Errorf("git commit failed: %w", err)
		}

		fmt.Printf("\nCommitted %d note(s) with message: %s\n", len(mdFiles), message)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
