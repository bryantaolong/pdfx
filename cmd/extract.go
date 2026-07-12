package cmd

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

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract specified pages from a PDF and merge them into a new file",
	RunE:  runExtract,
}

var extractName string
var extractPages string
var extractOutput string

func init() {
	extractCmd.Flags().StringVarP(&extractName, "name", "n", "", "Input PDF file path (required)")
	extractCmd.Flags().StringVarP(&extractPages, "pages", "p", "", "Comma-separated page numbers, e.g. 1,2,3,4 (required)")
	extractCmd.Flags().StringVarP(&extractOutput, "output", "o", "", "Output file path (default: <input>_extracted.pdf)")
	_ = extractCmd.MarkFlagRequired("name")
	_ = extractCmd.MarkFlagRequired("pages")
	rootCmd.AddCommand(extractCmd)
}

func runExtract(_ *cobra.Command, _ []string) error {
	extractName = ensurePDFExt(extractName)
	if _, err := os.Stat(extractName); os.IsNotExist(err) {
		return fmt.Errorf("file '%s' does not exist", extractName)
	}

	if extractOutput == "" {
		ext := filepath.Ext(extractName)
		stem := extractName[:len(extractName)-len(ext)]
		extractOutput = stem + "_extracted.pdf"
	}
	extractOutput = ensurePDFExt(extractOutput)

	parts := strings.Split(extractPages, ",")
	pages := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, err := strconv.Atoi(p); err != nil {
			return fmt.Errorf("invalid page number: %s", p)
		}
		pages = append(pages, p)
	}

	if len(pages) == 0 {
		return fmt.Errorf("no valid pages specified")
	}

	// Ensure output directory exists
	outDir := filepath.Dir(extractOutput)
	if outDir != "." && outDir != "" {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	fmt.Printf("Extracting pages to: %s\n", extractOutput)

	pageNrs := make([]int, 0, len(pages))
	for _, p := range pages {
		n, _ := strconv.Atoi(p)
		pageNrs = append(pageNrs, n)
	}

	ctx, err := api.ReadContextFile(extractName)
	if err != nil {
		return fmt.Errorf("failed to read PDF: %w", err)
	}

	newCtx, err := pdfcpu.ExtractPages(ctx, pageNrs, false)
	if err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	if err := api.WriteContextFile(newCtx, extractOutput); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	fmt.Printf("Extract complete! Saved to: %s\n", extractOutput)
	return nil
}
