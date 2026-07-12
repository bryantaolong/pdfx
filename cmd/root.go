package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pdfx",
	Short: "A Git-like CLI tool for merging, splitting, and extracting PDF pages",
}

func Execute() error {
	return rootCmd.Execute()
}

func ensurePDFExt(name string) string {
	if !strings.HasSuffix(strings.ToLower(name), ".pdf") {
		return name + ".pdf"
	}
	return name
}

func init() {
	// Global flags could be added here
}
