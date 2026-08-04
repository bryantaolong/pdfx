package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/cobra"
)

func NewCmdMerge() *cobra.Command {
	var mergeDir, mergeOutput string

	cmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge all PDF files in a directory into one file",
		RunE: func(_ *cobra.Command, _ []string) error {
			if mergeDir == "" {
				mergeDir = "."
			}
			if _, err := os.Stat(mergeDir); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("directory '%s' does not exist", mergeDir)
				}
				return fmt.Errorf("cannot access directory '%s': %w", mergeDir, err)
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
			mergeOutput = EnsurePDFExt(mergeOutput)

			// Ensure output directory exists
			outDir := filepath.Dir(mergeOutput)
			if outDir != "." && outDir != "" {
				if err := os.MkdirAll(outDir, 0755); err != nil {
					return fmt.Errorf("failed to create output directory: %w", err)
				}
			}

			// Check if output file already exists
			if _, err := os.Stat(mergeOutput); err == nil {
				return fmt.Errorf("output file '%s' already exists", mergeOutput)
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("cannot access output file '%s': %w", mergeOutput, err)
			}

			// Check if output conflicts with any input file
			absOutput, err := filepath.Abs(mergeOutput)
			if err != nil {
				return fmt.Errorf("failed to resolve output path: %w", err)
			}
			for _, m := range matches {
				absInput, err := filepath.Abs(m)
				if err != nil {
					return fmt.Errorf("failed to resolve input path '%s': %w", m, err)
				}
				if absInput == absOutput {
					return fmt.Errorf("output file '%s' conflicts with input file", mergeOutput)
				}
			}

			if err := api.MergeCreateFile(matches, mergeOutput, false, nil); err != nil {
				return fmt.Errorf("merge failed: %w", err)
			}

			fmt.Printf("\nMerge complete! Saved to: %s\n", mergeOutput)
			return nil
		},
	}

	cmd.Flags().StringVarP(&mergeDir, "dir", "d", "", "Directory containing PDF files (default: current directory)")
	cmd.Flags().StringVarP(&mergeOutput, "output", "o", "", "Output file path (default: merged.pdf)")
	return cmd
}
