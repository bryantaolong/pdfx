package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/cobra"
)

func NewCmdSplit() *cobra.Command {
	var splitFrom int

	cmd := &cobra.Command{
		Use:   "split <file.pdf>",
		Short: "Split a PDF into two files at a specified page number",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			splitName := args[0]
			splitName = EnsurePDFExt(splitName)
			if _, err := os.Stat(splitName); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("file '%s' does not exist", splitName)
				}
				return fmt.Errorf("cannot access file '%s': %w", splitName, err)
			}

			// Validate split page number.
			// A split at page 1 would produce an empty first file (pages 1..0),
			// which pdfcpu silently writes as a 0-page PDF, so reject it.
			if splitFrom < 2 {
				return fmt.Errorf("--from must be >= 2 (a split at page 1 would produce an empty first file), got %d", splitFrom)
			}

			dir := filepath.Dir(splitName)
			base := filepath.Base(splitName)
			stem := base[:len(base)-len(filepath.Ext(base))]
			if stem == "" {
				return fmt.Errorf("invalid PDF file name: %s", splitName)
			}

			totalPages, err := api.PageCountFile(splitName)
			if err != nil {
				return fmt.Errorf("failed to read PDF page count: %w", err)
			}
			if splitFrom > totalPages {
				return fmt.Errorf("--from must be <= total pages (%d), got %d", totalPages, splitFrom)
			}

			// Snapshot existing files so only newly generated ones are reported
			existing := map[string]bool{}
			prev, err := filepath.Glob(filepath.Join(dir, stem+"_*.pdf"))
			if err == nil {
				for _, p := range prev {
					existing[filepath.Base(p)] = true
				}
			}

			// Split at the specified page number
			if err := api.SplitByPageNrFile(splitName, dir, []int{splitFrom}, nil); err != nil {
				return fmt.Errorf("split failed: %w", err)
			}

			fmt.Println("Split complete!")
			fmt.Printf("  Output directory: %s\n", dir)

			// List newly generated files
			pattern := filepath.Join(dir, stem+"_*.pdf")
			matches, err := filepath.Glob(pattern)
			if err != nil {
				return fmt.Errorf("failed to list generated files: %w", err)
			}
			var generated []string
			for _, m := range matches {
				if !existing[filepath.Base(m)] {
					generated = append(generated, m)
				}
			}
			if len(generated) > 0 {
				sort.Strings(generated)
				fmt.Println("  Generated files:")
				for _, m := range generated {
					fmt.Printf("    - %s\n", filepath.Base(m))
				}
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&splitFrom, "from", "f", 0, "Start page of the second file, 1-based (required)")
	_ = cmd.MarkFlagRequired("from")
	return cmd
}
