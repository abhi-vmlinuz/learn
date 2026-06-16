package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"learn/internal/config"
	"learn/internal/editor"
	"learn/internal/fzf"
	"learn/internal/file"
	"learn/internal/template"

	"github.com/spf13/cobra"
)

var noEdit bool

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new note",
	Long:  "Create a new note using a template, with interactive category and template selection.",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireDeps("fzf", "EDITOR")

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		templatesDir := config.TemplatesDir()

		// List available templates and categories
		templates, err := template.ListAvailable(templatesDir)
		if err != nil {
			return fmt.Errorf("failed to list templates: %w", err)
		}
		if len(templates) == 0 {
			return fmt.Errorf("no templates found in %s", templatesDir)
		}

		categories := file.ListCategories(cfg.Repo.Root)
		if len(categories) == 0 {
			return fmt.Errorf("no categories found. Run 'learn init' first")
		}

		// Step 1: Select category
		category, err := fzf.Select(categories, "Select category")
		if err != nil {
			return err
		}

		// Step 2: Pick template — use matching one, fallback to general
		selectedTemplate := "general"
		for _, t := range templates {
			if t == category {
				selectedTemplate = t
				break
			}
		}

		// Step 3: Prompt for title
		fmt.Print("Note title: ")
		reader := bufio.NewReader(os.Stdin)
		title, _ := reader.ReadString('\n')
		title = strings.TrimSpace(title)
		if title == "" {
			return fmt.Errorf("title cannot be empty")
		}

		// Step 4: Load and render template (with fallback)
		tmplContent, err := template.LoadTemplate(templatesDir, selectedTemplate)
		if err != nil {
			// Fallback to general if template not found
			tmplContent, err = template.LoadTemplate(templatesDir, "general")
			if err != nil {
				return fmt.Errorf("no template found for %q and general template missing", selectedTemplate)
			}
			selectedTemplate = "general"
		}

		rendered := template.Render(tmplContent, title, category, selectedTemplate)

		// Step 5: Write file
		filename := makeFilename(title)
		categoryDir := filepath.Join(cfg.Repo.Root, category)
		filePath := filepath.Join(categoryDir, filename)

		if err := writeFile(filePath, rendered); err != nil {
			return fmt.Errorf("failed to write note: %w", err)
		}

		fmt.Printf("Created: %s\n", filePath)

		// Step 6: Open in editor
		if !noEdit {
			editor.OpenInEditor(filePath)
		}

		return nil
	},
}

func init() {
	newCmd.Flags().BoolVar(&noEdit, "no-edit", false, "Skip opening in editor")
	rootCmd.AddCommand(newCmd)
}

func makeFilename(title string) string {
	date := getCurrentDate()
	slug := slugify(title)
	return fmt.Sprintf("%s-%s.md", date, slug)
}

func writeFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}
