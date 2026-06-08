package cmd

import (
	"fmt"
	"path/filepath"

	"learn/internal/config"
	"learn/internal/fzf"
	"learn/internal/file"

	"github.com/spf13/cobra"
)

var recentCmd = &cobra.Command{
	Use:   "recent",
	Short: "Browse recently edited notes",
	Long:  "List recently modified notes sorted by time, with fzf selection and bat preview.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		learningDir := filepath.Join(cfg.Repo.Root, "learning")
		files, err := file.ListMarkdownFilesSorted(learningDir)
		if err != nil {
			return fmt.Errorf("failed to list notes: %w", err)
		}

		if len(files) == 0 {
			return fmt.Errorf("no notes found")
		}

		selected, err := fzf.SelectWithPreview(files, "Recent notes")
		if err != nil {
			return err
		}

		fmt.Printf("Opening: %s\n", selected)
		openInEditor(selected)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(recentCmd)
}
