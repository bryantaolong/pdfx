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

func NewCmdExtract() *cobra.Command {
	var extractName, extractPages, extractOutput string

	cmd := &cobra.Command{
		Use:   "extract",
		Short: "Extract specified pages from a PDF and merge them into a new file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			// 将包级 flag 值拷入局部变量，并重置包级变量，
			// 避免 cobra 在后续 Execute() 中复用上次残留值。
			name := extractName
			pages := extractPages
			output := extractOutput
			extractName = ""
			extractPages = ""
			extractOutput = ""

			if name == "" {
				return fmt.Errorf("required flag \"name\" not set")
			}
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
				output = stem + "_extracted.pdf"
			}
			output = EnsurePDFExt(output)

			// Check if output file already exists
			if _, err := os.Stat(output); err == nil {
				return fmt.Errorf("output file '%s' already exists", output)
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("cannot access output file '%s': %w", output, err)
			}

			// Check if output conflicts with the input file
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

			parts := strings.Split(pages, ",")
			pagesSlice := make([]string, 0, len(parts))
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
				pagesSlice = append(pagesSlice, p)
			}

			if len(pagesSlice) == 0 {
				return fmt.Errorf("no valid pages specified")
			}

			// Ensure output directory exists
			outDir := filepath.Dir(output)
			if outDir != "." && outDir != "" {
				if err := os.MkdirAll(outDir, 0755); err != nil {
					return fmt.Errorf("failed to create output directory: %w", err)
				}
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

			fmt.Printf("Extracting pages to: %s\n", output)

			pageNrs := make([]int, 0, len(pagesSlice))
			seen := make(map[int]struct{})
			for _, p := range pagesSlice {
				n, _ := strconv.Atoi(p)
				if _, ok := seen[n]; !ok {
					seen[n] = struct{}{}
					pageNrs = append(pageNrs, n)
				}
			}

			ctx, err := api.ReadContextFile(name)
			if err != nil {
				return fmt.Errorf("failed to read PDF: %w", err)
			}

			newCtx, err := pdfcpu.ExtractPages(ctx, pageNrs, false)
			if err != nil {
				return fmt.Errorf("extract failed: %w", err)
			}

			if err := api.WriteContextFile(newCtx, output); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}

			fmt.Printf("Extract complete! Saved to: %s\n", output)
			return nil
		},
	}

	cmd.Flags().StringVarP(&extractName, "name", "n", "", "Input PDF file path (required)")
	cmd.Flags().StringVarP(&extractPages, "pages", "p", "", "Comma-separated page numbers, e.g. 1,2,3,4 (required)")
	cmd.Flags().StringVarP(&extractOutput, "output", "o", "", "Output file path (default: <input>_extracted.pdf)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("pages")
	return cmd
}
