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
	var splitName string
	var splitFrom int

	cmd := &cobra.Command{
		Use:   "split",
		Short: "Split a PDF into two files at a specified page number",
		RunE: func(_ *cobra.Command, _ []string) error {
			splitName = EnsurePDFExt(splitName)
			if _, err := os.Stat(splitName); os.IsNotExist(err) {
				return fmt.Errorf("file '%s' does not exist", splitName)
			}

			dir := filepath.Dir(splitName)
			stem := filepath.Base(splitName[:len(splitName)-len(filepath.Ext(splitName))])

			// Split at the specified page number
			if err := api.SplitByPageNrFile(splitName, dir, []int{splitFrom}, nil); err != nil {
				return fmt.Errorf("split failed: %w", err)
			}

			fmt.Println("Split complete!")
			fmt.Printf("  Output directory: %s\n", dir)

			// List generated files
			matches, _ := filepath.Glob(filepath.Join(dir, stem+"_*.pdf"))
			if len(matches) > 0 {
				sort.Strings(matches)
				fmt.Println("  Generated files:")
				for _, m := range matches {
					fmt.Printf("    - %s\n", filepath.Base(m))
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&splitName, "name", "n", "", "Input PDF file path (required)")
	cmd.Flags().IntVarP(&splitFrom, "from", "f", 0, "Start page of the second file, 1-based (required)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("from")
	return cmd
}
