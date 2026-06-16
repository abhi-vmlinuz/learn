package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"learn/internal/config"
	"learn/internal/fzf"
	"learn/internal/file"

	"github.com/spf13/cobra"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:   "delete [filepath]",
	Short: "Delete a note",
	Long:  "Move a note to trash (or permanently delete with --force).",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireDeps("fzf")

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
			target, err = fzf.SelectWithPreview(files, "Delete note")
			if err != nil {
				return err
			}
		}

		if _, err := os.Stat(target); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", target)
		}

		// Confirm unless --force
		if !deleteForce {
			fmt.Printf("Delete %s? [y/N] ", filepath.Base(target))
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		if err := os.Remove(target); err != nil {
			return fmt.Errorf("failed to delete: %w", err)
		}

		fmt.Printf("Deleted: %s\n", filepath.Base(target))
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Skip confirmation")
	rootCmd.AddCommand(deleteCmd)
}
