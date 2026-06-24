package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"learn/internal/config"
	"learn/internal/editor"
	"learn/internal/fzf"
	"learn/internal/file"
	"learn/internal/git"

	"github.com/spf13/cobra"
)

var endCmd = &cobra.Command{
	Use:   "end",
	Short: "End your learning session",
	Long: `End-of-day workflow: reflect on what you learned, commit everything,
view your stats, and optionally shut down.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if !git.IsRepo(cfg.Repo.Root) {
			return fmt.Errorf("not a git repository: %s", cfg.Repo.Root)
		}

		fmt.Println("=== End of Learning Session ===")
		fmt.Println()

		today := time.Now().Format("2006-01-02")

		// Step 1: Prompt for reflection
		fmt.Println("What's on your mind today?")
		fmt.Println("(Anything goes — technical, personal, random thoughts. Press Enter twice to finish, or 'skip')")
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)
		var lines []string
		emptyCount := 0
		for {
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line == "skip" {
				break
			}
			if line == "" {
				emptyCount++
				if emptyCount >= 2 {
					break
				}
				lines = append(lines, "")
				continue
			}
			emptyCount = 0
			lines = append(lines, line)
		}

		// Step 2: Save journal entry if not skipped
		if len(lines) > 0 {
			journalDir := filepath.Join(cfg.Repo.Root, "journal")
			if err := os.MkdirAll(journalDir, 0755); err != nil {
				return fmt.Errorf("failed to create journal dir: %w", err)
			}

			now := time.Now()
			filename := fmt.Sprintf("%s.md", today)
			filePath := filepath.Join(journalDir, filename)

			// Check if entry already exists
			if _, err := os.Stat(filePath); err == nil {
				// Append to existing entry
				existing, _ := os.ReadFile(filePath)
				content := string(existing) + "\n\n---\n\n" + strings.Join(lines, "\n") + "\n"
				os.WriteFile(filePath, []byte(content), 0644)
				fmt.Printf("\nAppended to: %s\n", filePath)
			} else {
				// Create new entry
				content := fmt.Sprintf(`---
title: "Journal — %s"
date: %s
category: journal
created_at: %s
tags: ["journal", "daily"]
status: active
---

# Journal — %s

%s
`, today, today, now.Format(time.RFC3339), today, strings.Join(lines, "\n"))

				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					return fmt.Errorf("failed to write journal: %w", err)
				}
				fmt.Printf("\nJournal saved: %s\n", filePath)
			}
		}

		// Step 3: Commit everything
		fmt.Println()
		fmt.Println("Committing notes...")
		files, err := git.Status(cfg.Repo.Root)
		if err == nil && len(files) > 0 {
			var noteFiles []string
			for _, f := range files {
				if strings.HasSuffix(f, ".md") || strings.HasSuffix(f, ".pdf") {
					noteFiles = append(noteFiles, f)
				}
			}

			if len(noteFiles) > 0 {
				commitMsg := fmt.Sprintf("learn: %s session notes", today)
				if err := git.AddFiles(cfg.Repo.Root, noteFiles); err == nil {
					git.Commit(cfg.Repo.Root, commitMsg)
					fmt.Printf("Committed %d file(s): %s\n", len(noteFiles), commitMsg)
				}
			} else {
				fmt.Println("No notes to commit.")
			}
		} else {
			fmt.Println("No changes to commit.")
		}

		// Step 4: Show stats
		fmt.Println()
		allFiles, _ := file.ListMarkdownFiles(cfg.Repo.Root)
		fmt.Printf("Total notes: %d\n", len(allFiles))

		// Step 5: Shutdown prompt
		fmt.Println()
		if editor.IsTerminal() {
			choice, err := fzf.Select([]string{"no", "shutdown"}, "Power off?")
			if err == nil && choice == "shutdown" {
				fmt.Println("Shutting down in 10 seconds... (Ctrl+C to cancel)")
				time.Sleep(2 * time.Second)
				shutdown()
			}
		}

		fmt.Println()
		fmt.Println("Good session. See you tomorrow.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(endCmd)
}

func shutdown() {
	cmd := exec.Command("sudo", "shutdown", "-h", "+0")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
