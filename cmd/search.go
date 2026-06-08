package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"learn/internal/config"
	"learn/internal/fzf"

	"github.com/spf13/cobra"
)

var searchCategory string

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Full-text search notes",
	Long:  "Search notes using ripgrep with fzf selection and bat preview.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		// Build search path
		searchPath := filepath.Join(cfg.Repo.Root, "learning")
		if searchCategory != "" {
			searchPath = filepath.Join(searchPath, searchCategory)
			if _, err := os.Stat(searchPath); os.IsNotExist(err) {
				return fmt.Errorf("category %q not found", searchCategory)
			}
		}

		// Run ripgrep
		rgArgs := []string{"--files-with-matches"}
		if query != "" {
			rgArgs = append(rgArgs, query)
		}
		rgArgs = append(rgArgs, searchPath)

		rgCmd := exec.Command("rg", rgArgs...)
		rgOut, err := rgCmd.Output()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
				return fmt.Errorf("no matches found")
			}
			return fmt.Errorf("ripgrep failed: %w", err)
		}

		files := strings.Split(strings.TrimSpace(string(rgOut)), "\n")
		if len(files) == 0 || (len(files) == 1 && files[0] == "") {
			return fmt.Errorf("no matches found")
		}

		// Present via fzf with bat preview
		selected, err := fzf.SelectWithPreview(files, "Search results")
		if err != nil {
			return err
		}

		fmt.Printf("Opening: %s\n", selected)
		openInEditor(selected)

		return nil
	},
}

func init() {
	searchCmd.Flags().StringVar(&searchCategory, "category", "", "Filter search to specific category")
	rootCmd.AddCommand(searchCmd)
}
