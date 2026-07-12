package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge all PDF files in a directory into one file",
	RunE:  runMerge,
}

var mergeDir string
var mergeOutput string

func init() {
	mergeCmd.Flags().StringVarP(&mergeDir, "dir", "d", "", "Directory containing PDF files (default: current directory)")
	mergeCmd.Flags().StringVarP(&mergeOutput, "output", "o", "", "Output file path (default: merged.pdf)")
	rootCmd.AddCommand(mergeCmd)
}

func runMerge(_ *cobra.Command, _ []string) error {
	if mergeDir == "" {
		mergeDir = "."
	}
	if _, err := os.Stat(mergeDir); os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' does not exist", mergeDir)
	}

	pattern := filepath.Join(mergeDir, "*.pdf")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		fmt.Printf("Warning: no PDF files found in '%s'.\n", mergeDir)
		return nil
	}

	sort.Strings(matches)
	fmt.Printf("Found %d PDF files:\n", len(matches))
	for _, m := range matches {
		fmt.Printf("  - %s\n", filepath.Base(m))
	}

	if mergeOutput == "" {
		mergeOutput = "merged.pdf"
	}
	mergeOutput = ensurePDFExt(mergeOutput)

	if err := api.MergeCreateFile(matches, mergeOutput, false, nil); err != nil {
		return fmt.Errorf("merge failed: %w", err)
	}

	fmt.Printf("\nMerge complete! Saved to: %s\n", mergeOutput)
	return nil
}
