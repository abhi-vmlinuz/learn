package editor

import (
	"fmt"
	"os"
	"os/exec"
)

// GetEditor returns the user's preferred editor, or a fallback.
func GetEditor() string {
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	if HasBinary("nvim") {
		return "nvim"
	}
	if HasBinary("vim") {
		return "vim"
	}
	if HasBinary("vi") {
		return "vi"
	}
	return ""
}

// HasBinary checks if a binary is available in PATH.
func HasBinary(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// OpenInEditor opens a file in $EDITOR.
func OpenInEditor(path string) {
	editor := GetEditor()
	if editor == "" {
		fmt.Println("No editor found. Set $EDITOR environment variable.")
		return
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// OpenInViewer opens a markdown file in glow, falls back to $EDITOR.
func OpenInViewer(path string) {
	if HasBinary("glow") {
		cmd := exec.Command("glow", "-p", path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return
	}
	fmt.Println("glow not found, falling back to $EDITOR")
	OpenInEditor(path)
}

// OpenInPDFViewer opens a PDF file. Prefers tdf (terminal), falls back to xdg-open.
func OpenInPDFViewer(path string) {
	if HasBinary("tdf") {
		cmd := exec.Command("tdf", path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return
	}
	if HasBinary("xdg-open") {
		cmd := exec.Command("xdg-open", path)
		cmd.Start()
		return
	}
	fmt.Printf("No PDF viewer found. File saved at: %s\n", path)
}
