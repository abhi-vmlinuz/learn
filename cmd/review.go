package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"learn/internal/config"
	"learn/internal/fzf"
	"learn/internal/file"

	"github.com/spf13/cobra"
)

var reviewDays int

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review old notes for spaced repetition",
	Long:  "Randomly select notes older than N days for knowledge review.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		learningDir := cfg.Repo.Root
		allFiles, err := file.ListMarkdownFiles(learningDir)
		if err != nil {
			return fmt.Errorf("failed to list notes: %w", err)
		}

		cutoff := time.Now().AddDate(0, 0, -reviewDays)

		var candidates []string
		for _, f := range allFiles {
			info, err := os.Stat(f)
			if err != nil {
				continue
			}
			if info.ModTime().Before(cutoff) {
				candidates = append(candidates, f)
			}
		}

		if len(candidates) == 0 {
			return fmt.Errorf("no notes older than %d days found", reviewDays)
		}

		// Shuffle and take up to 20
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})
		if len(candidates) > 20 {
			candidates = candidates[:20]
		}

		selected, err := fzf.SelectWithPreview(candidates, fmt.Sprintf("Review (older than %d days)", reviewDays))
		if err != nil {
			return err
		}

		fmt.Printf("Opening: %s\n", selected)
		openInViewer(selected)

		return nil
	},
}

func init() {
	reviewCmd.Flags().IntVar(&reviewDays, "days", 7, "Minimum age in days for review candidates")
	rootCmd.AddCommand(reviewCmd)
}
