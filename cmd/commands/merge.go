package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

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

			// Match all entries and filter by ".pdf" case-insensitively,
			// so files like "DOC.PDF" are found on any platform.
			entries, err := filepath.Glob(filepath.Join(mergeDir, "*"))
			if err != nil {
				return err
			}
			matches := make([]string, 0, len(entries))
			for _, e := range entries {
				if strings.EqualFold(filepath.Ext(e), ".pdf") {
					matches = append(matches, e)
				}
			}

			if len(matches) == 0 {
				fmt.Printf("Warning: no PDF files found in '%s'.\n", mergeDir)
				return nil
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

			// Resolve output path once for input filtering
			absOutput, err := filepath.Abs(mergeOutput)
			if err != nil {
				return fmt.Errorf("failed to resolve output path: %w", err)
			}

			sort.Strings(matches)

			// Exclude the output file itself and any invalid (unreadable or 0-page) PDFs,
			// so a single broken file cannot abort the whole merge.
			inputs := make([]string, 0, len(matches))
			skipped := 0
			for _, m := range matches {
				absInput, err := filepath.Abs(m)
				if err != nil {
					return fmt.Errorf("failed to resolve input path '%s': %w", m, err)
				}
				if absInput == absOutput {
					continue
				}
				pages, err := api.PageCountFile(m)
				if err != nil || pages < 1 {
					skipped++
					fmt.Printf("Warning: skipping '%s' (not a valid PDF)\n", filepath.Base(m))
					continue
				}
				inputs = append(inputs, m)
			}

			if len(inputs) == 0 {
				return fmt.Errorf("no valid PDF files to merge in '%s'", mergeDir)
			}

			fmt.Printf("Found %d PDF files:\n", len(inputs))
			for _, m := range inputs {
				fmt.Printf("  - %s\n", filepath.Base(m))
			}
			if skipped > 0 {
				fmt.Printf("Warning: skipped %d invalid file(s).\n", skipped)
			}

			if err := api.MergeCreateFile(inputs, mergeOutput, false, nil); err != nil {
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
