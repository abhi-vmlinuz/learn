package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"learn/internal/config"
	"learn/internal/file"

	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show learning statistics",
	Long:  "Display note counts, category breakdown, and learning streaks.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		learningDir := filepath.Join(cfg.Repo.Root, "learning")
		allFiles, err := file.ListMarkdownFiles(learningDir)
		if err != nil {
			return fmt.Errorf("failed to list notes: %w", err)
		}

		if len(allFiles) == 0 {
			fmt.Println("No notes found.")
			return nil
		}

		// Count by category
		categoryCounts := make(map[string]int)
		var latestFile string
		var latestTime time.Time

		for _, f := range allFiles {
			rel, _ := filepath.Rel(learningDir, f)
			parts := strings.SplitN(rel, string(filepath.Separator), 2)
			category := "uncategorized"
			if len(parts) > 1 {
				category = parts[0]
			}
			categoryCounts[category]++

			info, err := os.Stat(f)
			if err == nil && info.ModTime().After(latestTime) {
				latestTime = info.ModTime()
				latestFile = filepath.Base(f)
			}
		}

		// Calculate streaks
		currentStreak, longestStreak := calculateStreaks(learningDir)

		// Print stats
		fmt.Printf("Total Notes: %d\n\n", len(allFiles))

		fmt.Println("Categories:")
		// Sort categories alphabetically
		categories := make([]string, 0, len(categoryCounts))
		for cat := range categoryCounts {
			categories = append(categories, cat)
		}
		sort.Strings(categories)
		for _, cat := range categories {
			fmt.Printf("  %-20s %d\n", cat, categoryCounts[cat])
		}

		fmt.Printf("\nCurrent Streak: %d days\n", currentStreak)
		fmt.Printf("Longest Streak: %d days\n", longestStreak)

		if latestFile != "" {
			fmt.Printf("\nLast Note:\n  %s\n", latestFile)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func calculateStreaks(learningDir string) (current, longest int) {
	dailyDir := filepath.Join(learningDir, "daily")
	entries, err := os.ReadDir(dailyDir)
	if err != nil {
		return 0, 0
	}

	var dates []time.Time
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		dateStr := strings.TrimSuffix(e.Name(), ".md")
		t, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			dates = append(dates, t)
		}
	}

	if len(dates) == 0 {
		return 0, 0
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	// Calculate longest streak
	longest = 1
	streak := 1
	for i := 1; i < len(dates); i++ {
		diff := dates[i].Sub(dates[i-1])
		if diff.Hours() <= 48 { // within 2 days (allowing for timezone)
			streak++
			if streak > longest {
				longest = streak
			}
		} else {
			streak = 1
		}
	}

	// Calculate current streak (from most recent date backwards)
	current = 1
	for i := len(dates) - 1; i > 0; i-- {
		diff := dates[i].Sub(dates[i-1])
		if diff.Hours() <= 48 {
			current++
		} else {
			break
		}
	}

	return current, longest
}
