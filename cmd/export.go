package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"learn/internal/config"
	"learn/internal/editor"
	"learn/internal/export"
	"learn/internal/fzf"
	"learn/internal/file"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [filepath]",
	Short: "Export a note to PDF",
	Long: `Export a markdown note to a beautifully styled PDF.

If no filepath is given, select a note interactively via fzf.
Requires wkhtmltopdf to be installed.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !export.IsAvailable() {
			return fmt.Errorf("wkhtmltopdf is not installed\nInstall it: sudo dnf install wkhtmltopdf")
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		var mdPath string

		if len(args) > 0 {
			mdPath = args[0]
			if !filepath.IsAbs(mdPath) {
				mdPath = filepath.Join(cfg.Repo.Root, mdPath)
			}
		} else {
			// Interactive selection
			files, err := file.ListMarkdownFilesSorted(cfg.Repo.Root)
			if err != nil {
				return fmt.Errorf("failed to list notes: %w", err)
			}
			if len(files) == 0 {
				return fmt.Errorf("no notes found")
			}

			mdPath, err = fzf.SelectWithPreview(files, "Export to PDF")
			if err != nil {
				return err
			}
		}

		// Check file exists
		if _, err := os.Stat(mdPath); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", mdPath)
		}

		// Generate output path
		baseName := strings.TrimSuffix(filepath.Base(mdPath), ".md")
		outDir := filepath.Dir(mdPath)
		outPath := filepath.Join(outDir, baseName+".pdf")

		fmt.Printf("Exporting: %s\n", mdPath)
		fmt.Printf("     To:   %s\n", outPath)

		if err := export.GeneratePDF(mdPath, outPath); err != nil {
			return fmt.Errorf("export failed: %w", err)
		}

		fmt.Println("Done.")

		// Offer to open the PDF (only in interactive mode)
		if isTerminal() {
			openChoices := []string{"no"}
			if editor.HasBinary("tdf") {
				openChoices = append([]string{"tdf (terminal)"}, openChoices...)
			}
			openChoices = append([]string{"default app"}, openChoices...)

			if len(openChoices) > 1 {
				choice, err := fzf.Select(openChoices, "Open PDF")
				if err == nil {
					switch choice {
					case "tdf (terminal)":
						editor.OpenInPDFViewer(outPath)
					case "default app":
						editor.OpenInPDFViewer(outPath)
					}
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func isTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
