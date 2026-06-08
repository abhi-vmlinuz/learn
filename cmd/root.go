package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "learn",
	Short: "Git-first engineering knowledge base CLI",
	Long:  "Learn is a lightweight, Git-first knowledge base CLI for DevOps, platform engineers, and Linux system administrators.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
