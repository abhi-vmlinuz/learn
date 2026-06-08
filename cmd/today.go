package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"learn/internal/config"
	"learn/internal/template"

	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Create today's daily journal entry",
	Long:  "Create or open today's daily journal entry.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		templatesDir := config.TemplatesDir()
		date := time.Now().Format("2006-01-02")
		filename := date + ".md"
		filePath := filepath.Join(cfg.Repo.Root, "daily", filename)

		// Check if file already exists
		if _, err := os.Stat(filePath); err == nil {
			fmt.Printf("Journal entry for %s already exists.\n", date)
			fmt.Printf("Opening: %s\n", filePath)
			openInEditor(filePath)
			return nil
		}

		// Load daily template
		tmplContent, err := template.LoadTemplate(templatesDir, "daily")
		if err != nil {
			return fmt.Errorf("daily template not found: %w", err)
		}

		rendered := template.Render(tmplContent, "Daily Journal — "+date, "daily", "daily")

		if err := writeFile(filePath, rendered); err != nil {
			return fmt.Errorf("failed to write journal: %w", err)
		}

		fmt.Printf("Created: %s\n", filePath)
		openInEditor(filePath)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(todayCmd)
}

// getCurrentDate returns today's date as YYYY-MM-DD.
func getCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// slugify converts a title to a filename-safe slug.
func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
