package cmd

import (
	"fmt"

	"learn/internal/config"
	"learn/internal/editor"
	"learn/internal/fzf"
	"learn/internal/file"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [query]",
	Short: "Edit an existing note",
	Long:  "Browse and open a note in $EDITOR. Optionally filter with a search query.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireDeps("fzf", "EDITOR")

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		files, err := file.ListMarkdownFilesSorted(cfg.Repo.Root)
		if err != nil {
			return fmt.Errorf("failed to list notes: %w", err)
		}

		if len(files) == 0 {
			return fmt.Errorf("no notes found")
		}

		selected, err := fzf.SelectWithPreview(files, "Edit note")
		if err != nil {
			return err
		}

		fmt.Printf("Editing: %s\n", selected)
		editor.OpenInEditor(selected)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
