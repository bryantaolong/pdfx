package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/spf13/cobra"
)

func NewCmdRemove() *cobra.Command {
	var removePages, removeOutput string

	cmd := &cobra.Command{
		Use:   "remove <file.pdf>",
		Short: "Remove specified pages from a PDF and save the result to a new file",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			pages := removePages
			output := removeOutput
			removePages = ""
			removeOutput = ""

			if pages == "" {
				return fmt.Errorf("required flag \"pages\" not set")
			}

			name = EnsurePDFExt(name)
			if _, err := os.Stat(name); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("file '%s' does not exist", name)
				}
				return fmt.Errorf("cannot access file '%s': %w", name, err)
			}

			if output == "" {
				ext := filepath.Ext(name)
				stem := name[:len(name)-len(ext)]
				if stem == "" {
					return fmt.Errorf("invalid PDF file name: %s", name)
				}
				output = stem + "_removed.pdf"
			}
			output = EnsurePDFExt(output)

			// Check if output conflicts with the input file before auto-incrementing
			absInput, err := filepath.Abs(name)
			if err != nil {
				return fmt.Errorf("failed to resolve input path: %w", err)
			}
			absOutput, err := filepath.Abs(output)
			if err != nil {
				return fmt.Errorf("failed to resolve output path: %w", err)
			}
			if absInput == absOutput {
				return fmt.Errorf("output file '%s' conflicts with input file", output)
			}

			// Ensure output directory exists
			outDir := filepath.Dir(output)
			if outDir != "." && outDir != "" {
				if err := os.MkdirAll(outDir, 0755); err != nil {
					return fmt.Errorf("failed to create output directory: %w", err)
				}
			}

			// Check if output file already exists and auto-increment suffix if needed
			baseOutput := output
			counter := 1
			for {
				if _, err := os.Stat(output); err != nil {
					if os.IsNotExist(err) {
						break
					}
					return fmt.Errorf("cannot access output file '%s': %w", output, err)
				}
				ext := filepath.Ext(baseOutput)
				stem := baseOutput[:len(baseOutput)-len(ext)]
				output = fmt.Sprintf("%s_%d%s", stem, counter, ext)
				counter++
			}

			parts := strings.Split(pages, ",")
			pagesSlice := make([]string, 0, len(parts))
			seen := make(map[int]struct{})
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				n, err := strconv.Atoi(p)
				if err != nil {
					return fmt.Errorf("invalid page number: %s", p)
				}
				if n < 1 {
					return fmt.Errorf("invalid page number: %s (must be >= 1)", p)
				}
				if _, ok := seen[n]; !ok {
					seen[n] = struct{}{}
					pagesSlice = append(pagesSlice, p)
				}
			}

			if len(pagesSlice) == 0 {
				return fmt.Errorf("no valid pages specified")
			}

			// Validate page range
			totalPages, err := api.PageCountFile(name)
			if err != nil {
				return fmt.Errorf("failed to read PDF page count: %w", err)
			}
			for _, p := range pagesSlice {
				n, _ := strconv.Atoi(p)
				if n > totalPages {
					return fmt.Errorf("page number %d exceeds total pages (%d)", n, totalPages)
				}
			}

			fmt.Printf("Removing pages, output: %s\n", output)

			// Build set of pages to remove for fast lookup
			removeSet := make(map[int]struct{}, len(pagesSlice))
			for _, p := range pagesSlice {
				n, _ := strconv.Atoi(p)
				removeSet[n] = struct{}{}
			}

			// Determine pages to keep
			keepNrs := make([]int, 0, totalPages-len(pagesSlice))
			for i := 1; i <= totalPages; i++ {
				if _, ok := removeSet[i]; !ok {
					keepNrs = append(keepNrs, i)
				}
			}

			if len(keepNrs) == 0 {
				return fmt.Errorf("removing all pages would result in an empty PDF")
			}

			ctx, err := api.ReadContextFile(name)
			if err != nil {
				return fmt.Errorf("failed to read PDF: %w", err)
			}

			newCtx, err := pdfcpu.ExtractPages(ctx, keepNrs, false)
			if err != nil {
				return fmt.Errorf("remove failed: %w", err)
			}

			if err := api.WriteContextFile(newCtx, output); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}

			fmt.Printf("Remove complete! Saved to: %s\n", output)
			return nil
		},
	}

	cmd.Flags().StringVarP(&removePages, "pages", "p", "", "Comma-separated page numbers to remove, e.g. 1,2,3,4 (required)")
	cmd.Flags().StringVarP(&removeOutput, "output", "o", "", "Output file path (default: <input>_removed.pdf)")
	_ = cmd.MarkFlagRequired("pages")
	return cmd
}
