package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"learn/internal/config"
	"learn/internal/fzf"
	"learn/internal/template"

	"github.com/spf13/cobra"
)

var noEdit bool

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new note",
	Long:  "Create a new note using a template, with interactive category and template selection.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		templatesDir := config.TemplatesDir()

		// List templates
		templates, err := template.ListAvailable(templatesDir)
		if err != nil {
			return fmt.Errorf("failed to list templates: %w", err)
		}
		if len(templates) == 0 {
			return fmt.Errorf("no templates found in %s", templatesDir)
		}

		// Select template via fzf
		selectedTemplate, err := fzf.Select(templates, "Select template")
		if err != nil {
			return err
		}

		// Prompt for title
		fmt.Print("Note title: ")
		var title string
		fmt.Scanln(&title)
		if title == "" {
			return fmt.Errorf("title cannot be empty")
		}

		// Select category (mandatory)
		categories := []string{
			"linux", "aws", "docker", "kubernetes",
			"networking", "ctf", "troubleshooting", "challenge",
		}
		category, err := fzf.Select(categories, "Select category")
		if err != nil {
			return err
		}

		// Load and render template
		tmplContent, err := template.LoadTemplate(templatesDir, selectedTemplate)
		if err != nil {
			return err
		}

		rendered := template.Render(tmplContent, title, category, selectedTemplate)

		// Write file
		filename := makeFilename(title)
		categoryDir := filepath.Join(cfg.Repo.Root, category)
		filePath := filepath.Join(categoryDir, filename)

		if err := writeFile(filePath, rendered); err != nil {
			return fmt.Errorf("failed to write note: %w", err)
		}

		fmt.Printf("Created: %s\n", filePath)

		// Open in editor
		if !noEdit {
			openInEditor(filePath)
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

func openInEditor(path string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
