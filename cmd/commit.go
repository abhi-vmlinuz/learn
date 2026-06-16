package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"learn/internal/config"
	"learn/internal/git"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit [message]",
	Short: "Git add and commit notes",
	Long: `Stage all modified notes and commit with a message.

If no message is provided, you will be prompted to enter one.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireDeps("git")

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		repoRoot := cfg.Repo.Root

		// Check if it's actually a git repo
		if !git.IsRepo(repoRoot) {
			return fmt.Errorf("not a git repository: %s\nRun 'git init' in your learn directory first", repoRoot)
		}

		// Get modified/untracked files
		files, err := git.Status(repoRoot)
		if err != nil {
			return fmt.Errorf("failed to check git status: %w", err)
		}

		// Filter to markdown and pdf files only
		var noteFiles []string
		for _, f := range files {
			if strings.HasSuffix(f, ".md") || strings.HasSuffix(f, ".pdf") {
				noteFiles = append(noteFiles, f)
			}
		}

		if len(noteFiles) == 0 {
			fmt.Println("No modified or untracked notes to commit.")
			return nil
		}

		// Display summary
		fmt.Println("Will commit:")
		for _, f := range noteFiles {
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
			message = strings.Join(args, " ")
		} else {
			fmt.Print("Commit message: ")
			reader := bufio.NewReader(os.Stdin)
			message, _ = reader.ReadString('\n')
			message = strings.TrimSpace(message)
			if message == "" {
				return fmt.Errorf("commit message cannot be empty")
			}
		}

		// Stage only the markdown and pdf files and commit
		if err := git.AddFiles(repoRoot, noteFiles); err != nil {
			return fmt.Errorf("git add failed: %w", err)
		}

		if err := git.Commit(repoRoot, message); err != nil {
			return fmt.Errorf("git commit failed: %w", err)
		}

		fmt.Printf("\nCommitted %d note(s): %s\n", len(noteFiles), message)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
