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

var tagCmd = &cobra.Command{
	Use:   "tag [filepath]",
	Short: "Edit tags on a note",
	Long:  "Add or remove tags from an existing note's frontmatter.",
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
			target, err = fzf.SelectWithPreview(files, "Edit tags")
			if err != nil {
				return err
			}
		}

		data, err := os.ReadFile(target)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		content := string(data)

		// Parse current tags
		currentTags := extractTags(content)

		fmt.Printf("File: %s\n", filepath.Base(target))
		fmt.Printf("Current tags: %s\n", formatTags(currentTags))
		fmt.Println()
		fmt.Println("Enter new tags (comma-separated, or press Enter to keep current):")
		fmt.Print("> ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			fmt.Println("No changes.")
			return nil
		}

		// Parse new tags
		var newTags []string
		for _, t := range strings.Split(input, ",") {
			t = strings.TrimSpace(t)
			t = strings.Trim(t, "\"")
			if t != "" {
				newTags = append(newTags, t)
			}
		}

		// Replace tags in frontmatter
		updated := replaceTags(content, newTags)

		if err := os.WriteFile(target, []byte(updated), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Printf("Tags updated: %s\n", formatTags(newTags))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}

func extractTags(content string) []string {
	lines := strings.Split(content, "\n")
	inFrontmatter := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if inFrontmatter {
				break
			}
			inFrontmatter = true
			continue
		}
		if inFrontmatter && strings.HasPrefix(trimmed, "tags:") {
			val := strings.TrimSpace(strings.TrimPrefix(trimmed, "tags:"))
			val = strings.Trim(val, "[]")
			if val == "" {
				return nil
			}
			var tags []string
			for _, t := range strings.Split(val, ",") {
				t = strings.TrimSpace(t)
				t = strings.Trim(t, "\"")
				if t != "" {
					tags = append(tags, t)
				}
			}
			return tags
		}
	}
	return nil
}

func replaceTags(content string, tags []string) string {
	quoted := make([]string, len(tags))
	for i, t := range tags {
		quoted[i] = "\"" + t + "\""
	}
	newTagLine := "tags: [" + strings.Join(quoted, ", ") + "]"

	lines := strings.Split(content, "\n")
	inFrontmatter := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if inFrontmatter {
				break
			}
			inFrontmatter = true
			continue
		}
		if inFrontmatter && strings.HasPrefix(trimmed, "tags:") {
			lines[i] = newTagLine
			break
		}
	}
	return strings.Join(lines, "\n")
}

func formatTags(tags []string) string {
	if len(tags) == 0 {
		return "[]"
	}
	quoted := make([]string, len(tags))
	for i, t := range tags {
		quoted[i] = "\"" + t + "\""
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}
