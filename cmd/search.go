package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"learn/internal/config"
	"learn/internal/editor"
	"learn/internal/fzf"
	"learn/internal/file"

	"github.com/spf13/cobra"
)

var searchCategory string

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Full-text search notes",
	Long: `Search notes using ripgrep with fzf selection and bat preview.

Without a query, lists all notes for interactive browsing.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireDeps("fzf", "rg", "bat", "glow")

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		// Build search path
		searchPath := cfg.Repo.Root
		if searchCategory != "" {
			searchPath = filepath.Join(searchPath, searchCategory)
			if _, err := os.Stat(searchPath); os.IsNotExist(err) {
				return fmt.Errorf("category %q not found", searchCategory)
			}
		}

		var files []string

		if query == "" {
			// No query: list all markdown files
			files, err = file.ListMarkdownFiles(searchPath)
			if err != nil {
				return fmt.Errorf("failed to list notes: %w", err)
			}
		} else {
			// With query: use ripgrep for full-text search
			rgArgs := []string{"--files-with-matches", query, searchPath}
			rgCmd := exec.Command("rg", rgArgs...)
			rgOut, rgErr := rgCmd.Output()
			if rgErr != nil {
				if exitErr, ok := rgErr.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
					return fmt.Errorf("no matches found for %q", query)
				}
				return fmt.Errorf("ripgrep failed: %w", rgErr)
			}
			files = strings.Split(strings.TrimSpace(string(rgOut)), "\n")
		}

		if len(files) == 0 || (len(files) == 1 && files[0] == "") {
			return fmt.Errorf("no notes found")
		}

		// Present via fzf with bat preview
		selected, err := fzf.SelectWithPreview(files, "Notes")
		if err != nil {
			return err
		}

		// Ask read or edit
		action, err := fzf.Select([]string{"read", "edit"}, "Open mode")
		if err != nil {
			return err
		}

		if action == "read" {
			editor.OpenInViewer(selected)
		} else {
			editor.OpenInEditor(selected)
		}

		return nil
	},
}

func init() {
	searchCmd.Flags().StringVar(&searchCategory, "category", "", "Filter search to specific category")
	rootCmd.AddCommand(searchCmd)
}
