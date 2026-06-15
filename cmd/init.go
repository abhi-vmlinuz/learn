package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"learn/internal/config"
	"learn/internal/git"
	"learn/internal/template"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a learn repository",
	Long:  "Set up the learn directory structure, config, and templates in the current directory.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		// Check for .git
		if git.IsRepo(cwd) {
			fmt.Println("Git repo detected. Continuing...")
		} else {
			fmt.Println("WARNING: No .git directory found.")
			fmt.Println("To initialize a git repo, run:")
			fmt.Println("  git init && git remote add origin <url>")
			fmt.Println()
		}

		// Create category directories directly in current directory
		categories := []string{
			"aws",
			"linux",
			"docker",
			"kubernetes",
			"networking",
			"ctf",
			"troubleshooting",
			"daily",
			"challenge",
			"general",
		}

		for _, cat := range categories {
			dir := filepath.Join(cwd, cat)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create %s: %w", cat, err)
			}
		}

		// Save config
		cfg := &config.Config{
			Repo: config.RepoConfig{Root: cwd},
		}
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		// Copy bundled templates (backup existing first)
		templatesDir := config.TemplatesDir()
		if _, err := os.Stat(templatesDir); err == nil {
			// Templates dir exists, back it up
			backupDir := templatesDir + ".bak"
			if err := os.Rename(templatesDir, backupDir); err != nil {
				return fmt.Errorf("failed to backup templates: %w", err)
			}
			fmt.Printf("Existing templates backed up to: %s\n", backupDir)
		}
		if err := template.CopyDefaults(templatesDir); err != nil {
			return fmt.Errorf("failed to copy templates: %w", err)
		}

		fmt.Println("Learn repository initialized successfully.")
		fmt.Println()
		fmt.Printf("  Config:    %s\n", config.ConfigPath())
		fmt.Printf("  Templates: %s\n", templatesDir)
		fmt.Printf("  Root:      %s\n", cwd)
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  learn new       — create your first note")
		fmt.Println("  learn today     — start a daily journal")
		fmt.Println("  learn doctor    — check your environment")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
