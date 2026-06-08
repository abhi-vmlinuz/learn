package fzf

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Select presents items via fzf and returns the selected item.
func Select(items []string, prompt string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to select from")
	}

	input := strings.Join(items, "\n")
	cmd := exec.Command("fzf", "--prompt", prompt+"> ")
	cmd.Stdin = strings.NewReader(input)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return "", fmt.Errorf("selection cancelled")
		}
		return "", fmt.Errorf("fzf failed: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}

// SelectWithPreview presents items via fzf with bat preview.
func SelectWithPreview(items []string, prompt string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to select from")
	}

	input := strings.Join(items, "\n")
	cmd := exec.Command("fzf", "--prompt", prompt+"> ", "--preview", "bat --color=always {}")
	cmd.Stdin = strings.NewReader(input)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return "", fmt.Errorf("selection cancelled")
		}
		return "", fmt.Errorf("fzf failed: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}

// SelectWithPreviewRaw presents items via fzf with a custom preview command.
func SelectWithPreviewRaw(items []string, prompt, previewCmd string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to select from")
	}

	input := strings.Join(items, "\n")
	cmd := exec.Command("fzf", "--prompt", prompt+"> ", "--preview", previewCmd)
	cmd.Stdin = strings.NewReader(input)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return "", fmt.Errorf("selection cancelled")
		}
		return "", fmt.Errorf("fzf failed: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}

// IsAvailable checks if fzf is installed.
func IsAvailable() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}
