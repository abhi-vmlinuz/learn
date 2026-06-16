package cmd

import (
	"fmt"

	"learn/internal/config"
	"learn/internal/editor"
	"learn/internal/fzf"
	"learn/internal/file"

	"github.com/spf13/cobra"
)

var recentCmd = &cobra.Command{
	Use:   "recent",
	Short: "Browse recently edited notes",
	Long:  "List recently modified notes sorted by time, with fzf selection and bat preview.",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireDeps("fzf", "bat", "glow")

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		learningDir := cfg.Repo.Root
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
		editor.OpenInViewer(selected)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(recentCmd)
}
