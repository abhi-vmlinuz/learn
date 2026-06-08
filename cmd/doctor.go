package cmd

import (
	"fmt"
	"os/exec"

	"learn/internal/config"
	"learn/internal/fzf"
	"learn/internal/git"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check environment health",
	Long:  "Verify all dependencies and configuration are properly set up.",
	RunE: func(cmd *cobra.Command, args []string) error {
		allGood := true

		checks := []struct {
			name string
			fn   func() bool
		}{
			{"git", func() bool { return git.IsAvailable() }},
			{"fzf", func() bool { return fzf.IsAvailable() }},
			{"rg", func() bool { _, err := exec.LookPath("rg"); return err == nil }},
			{"bat", func() bool { _, err := exec.LookPath("bat"); return err == nil }},
			{"EDITOR", func() bool { _, err := exec.LookPath(getEditor()); return err == nil }},
		}

		for _, c := range checks {
			if c.fn() {
				fmt.Printf("  ✓ %s\n", c.name)
			} else {
				fmt.Printf("  ✗ %s\n", c.name)
				allGood = false
			}
		}

		// Check config
		_, err := config.Load()
		if err == nil {
			fmt.Println("  ✓ repository")
			fmt.Println("  ✓ config file")
		} else {
			fmt.Println("  ✗ repository (run 'learn init')")
			fmt.Println("  ✗ config file")
			allGood = false
		}

		fmt.Println()
		if allGood {
			fmt.Println("Repository healthy")
		} else {
			fmt.Println("Some checks failed. Fix the issues above.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func getEditor() string {
	editor := exec.Command("bash", "-c", "echo $EDITOR")
	out, err := editor.Output()
	if err != nil {
		return "vi"
	}
	result := string(out)
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	if result == "" {
		return "vi"
	}
	return result
}
