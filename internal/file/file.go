package file

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// TitleFromFilename extracts the title from a YYYY-MM-DD-title.md filename.
func TitleFromFilename(filename string) string {
	base := filepath.Base(filename)
	base = strings.TrimSuffix(base, ".md")

	// Match YYYY-MM-DD-title pattern
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-(.+)$`)
	matches := re.FindStringSubmatch(base)
	if len(matches) > 1 {
		return strings.ReplaceAll(matches[1], "-", " ")
	}
	return base
}

// DateFromFilename extracts the date from a YYYY-MM-DD-title.md filename.
func DateFromFilename(filename string) (time.Time, error) {
	base := filepath.Base(filename)
	re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})-`)
	matches := re.FindStringSubmatch(base)
	if len(matches) > 1 {
		return time.Parse("2006-01-02", matches[1])
	}
	return time.Time{}, fmt.Errorf("no date found in filename: %s", filename)
}

// MakeFilename creates a YYYY-MM-DD-title.md filename.
func MakeFilename(title string) string {
	date := time.Now().Format("2006-01-02")
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return fmt.Sprintf("%s-%s.md", date, slug)
}

// EnsureDir creates a directory if it doesn't exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// WriteFile writes content to a file, creating parent dirs if needed.
func WriteFile(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// ListCategories returns subdirectories of the repo root that are category dirs.
// Excludes hidden dirs, daily, and common non-category dirs.
func ListCategories(repoRoot string) []string {
	entries, err := os.ReadDir(repoRoot)
	if err != nil {
		return nil
	}

	skip := map[string]bool{
		".git": true,
	}

	var dirs []string
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") && !skip[e.Name()] {
			dirs = append(dirs, e.Name())
		}
	}
	return dirs
}

// ListMarkdownFiles returns all .md files under a directory recursively.
func ListMarkdownFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip errors
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// ListMarkdownFilesSorted returns .md files sorted by modification time (newest first).
func ListMarkdownFilesSorted(root string) ([]string, error) {
	files, err := ListMarkdownFiles(root)
	if err != nil {
		return nil, err
	}

	type fileWithTime struct {
		path    string
		modTime time.Time
	}

	var sorted []fileWithTime
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			continue
		}
		sorted = append(sorted, fileWithTime{f, info.ModTime()})
	}

	// Sort by mod time descending
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].modTime.After(sorted[i].modTime) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	result := make([]string, len(sorted))
	for i, s := range sorted {
		result[i] = s.path
	}
	return result, nil
}
