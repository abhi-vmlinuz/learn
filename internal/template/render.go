package template

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed defaults/*
var bundledTemplates embed.FS

// ListAvailable returns template names from the user's templates dir.
func ListAvailable(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates dir: %w", err)
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			names = append(names, strings.TrimSuffix(e.Name(), ".md"))
		}
	}
	return names, nil
}

// LoadTemplate reads a template file from the user's templates dir.
func LoadTemplate(dir, name string) (string, error) {
	path := filepath.Join(dir, name+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("template %q not found: %w", name, err)
	}
	return string(data), nil
}

// categoryTags maps categories to relevant auto-generated tags.
var categoryTags = map[string][]string{
	"linux":           {"linux", "sysadmin", "cli"},
	"aws":             {"aws", "cloud", "iaas"},
	"docker":          {"docker", "containers", "devops"},
	"kubernetes":      {"kubernetes", "k8s", "containers", "devops"},
	"networking":      {"networking", "tcp", "sysadmin"},
	"ctf":             {"ctf", "security", "challenge"},
	"troubleshooting": {"troubleshooting", "debugging"},
	"daily":           {"daily", "journal"},
	"challenge":       {"challenge", "learning"},
	"general":         {"general", "notes"},
}

// autoTags returns a YAML-formatted tag list for the given category and template.
func autoTags(category, templateName string) string {
	tags := []string{category}
	if templateName != category {
		tags = append(tags, templateName)
	}
	if extra, ok := categoryTags[category]; ok {
		for _, t := range extra {
			found := false
			for _, existing := range tags {
				if existing == t {
					found = true
					break
				}
			}
			if !found {
				tags = append(tags, t)
			}
		}
	}
	// Format as YAML list
	quoted := make([]string, len(tags))
	for i, t := range tags {
		quoted[i] = "\"" + t + "\""
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

// Render replaces placeholders in a template string.
func Render(content, title, category, templateName string) string {
	now := time.Now()
	tags := autoTags(category, templateName)
	replacements := map[string]string{
		"{title}":    title,
		"{date}":     now.Format("2006-01-02"),
		"{datetime}": now.Format(time.RFC3339),
		"{category}": category,
		"{tags}":     tags,
	}

	result := content
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// Fallback: if template still has old empty tags, inject auto-generated ones
	result = strings.ReplaceAll(result, "tags: []", "tags: "+tags)

	return result
}

// CopyDefaults copies bundled default templates to the user's templates dir.
func CopyDefaults(userDir string) error {
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return fmt.Errorf("failed to create templates dir: %w", err)
	}

	entries, err := bundledTemplates.ReadDir("defaults")
	if err != nil {
		return fmt.Errorf("failed to read bundled templates: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		data, err := bundledTemplates.ReadFile("defaults/" + entry.Name())
		if err != nil {
			return fmt.Errorf("failed to read bundled template %s: %w", entry.Name(), err)
		}

		dest := filepath.Join(userDir, entry.Name())
		if err := os.WriteFile(dest, data, 0644); err != nil {
			return fmt.Errorf("failed to write template %s: %w", entry.Name(), err)
		}
	}

	return nil
}
