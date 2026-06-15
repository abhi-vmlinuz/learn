package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"learn/internal/config"
	"learn/internal/fzf"
	"learn/internal/file"

	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move [filepath]",
	Short: "Move a note to a different category",
	Long:  "Move a note from one category to another.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		var target string

		if len(args) > 0 {
			target = args[0]
			if !filepath.IsAbs(target) {
				target = filepath.Join(cfg.Repo.Root, target)
			}
		} else {
			files, err := file.ListMarkdownFilesSorted(cfg.Repo.Root)
			if err != nil {
				return fmt.Errorf("failed to list notes: %w", err)
			}
			if len(files) == 0 {
				return fmt.Errorf("no notes found")
			}
			target, err = fzf.SelectWithPreview(files, "Move note")
			if err != nil {
				return err
			}
		}

		if _, err := os.Stat(target); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", target)
		}

		// Select destination category
		categories := file.ListCategories(cfg.Repo.Root)
		if len(categories) == 0 {
			return fmt.Errorf("no categories found")
		}

		dest, err := fzf.Select(categories, "Move to category")
		if err != nil {
			return err
		}

		// Check if already in that category
		rel, _ := filepath.Rel(cfg.Repo.Root, target)
		currentCat := filepath.Dir(rel)
		if currentCat == dest {
			fmt.Println("Already in that category.")
			return nil
		}

		// Move file
		destDir := filepath.Join(cfg.Repo.Root, dest)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create category dir: %w", err)
		}

		newPath := filepath.Join(destDir, filepath.Base(target))
		if err := os.Rename(target, newPath); err != nil {
			return fmt.Errorf("failed to move: %w", err)
		}

		fmt.Printf("Moved: %s -> %s/\n", filepath.Base(target), dest)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
}
