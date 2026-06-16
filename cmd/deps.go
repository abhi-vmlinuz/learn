package cmd

import (
	"fmt"
	"os"

	"learn/internal/editor"
	"learn/internal/fzf"
	"learn/internal/git"
)

// requireDeps checks that required tools are available, exits with clear message if not.
func requireDeps(tools ...string) {
	missing := []string{}
	for _, tool := range tools {
		switch tool {
		case "git":
			if !git.IsAvailable() {
				missing = append(missing, "git")
			}
		case "fzf":
			if !fzf.IsAvailable() {
				missing = append(missing, "fzf")
			}
		case "rg":
			if !editor.HasBinary("rg") {
				missing = append(missing, "rg (ripgrep)")
			}
		case "bat":
			if !editor.HasBinary("bat") {
				missing = append(missing, "bat")
			}
		case "glow":
			if !editor.HasBinary("glow") {
				missing = append(missing, "glow")
			}
		case "wkhtmltopdf":
			if !editor.HasBinary("wkhtmltopdf") {
				missing = append(missing, "wkhtmltopdf")
			}
		case "EDITOR":
			if editor.GetEditor() == "" {
				missing = append(missing, "$EDITOR")
			}
		}
	}

	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "Missing required dependencies: ")
		for i, m := range missing {
			if i > 0 {
				fmt.Fprintf(os.Stderr, ", ")
			}
			fmt.Fprintf(os.Stderr, "%s", m)
		}
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Run 'learn doctor' for install instructions.\n")
		os.Exit(1)
	}
}
