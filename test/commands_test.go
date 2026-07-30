package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bryantaolong/pdfx/cmd"
)

var projectRoot string

func TestMain(m *testing.M) {
	// Determine source PDFs directory before changing working directory
	// go test runs from the package directory, so we need to go up to project root
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	projectRoot = filepath.Join(cwd, "..")
	srcDir := filepath.Join(projectRoot, "pdfs")

	// Create a temp working directory for each test run
	dir, err := os.MkdirTemp("", "pdfx-test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	if err := os.Chdir(dir); err != nil {
		panic(err)
	}

	// Prepare test PDFs by copying from repo pdfs/
	for _, name := range []string{"a.pdf", "b.pdf", "c.pdf", "test_file.pdf"} {
		src := filepath.Join(srcDir, name)
		dst := filepath.Join(dir, name)
		data, err := os.ReadFile(src)
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			panic(err)
		}
	}

	// Run tests from the temp directory
	os.Exit(m.Run())
}

func runCLI(args ...string) error {
	os.Args = append([]string{"pdfx"}, args...)
	return cmd.Execute()
}

func TestVersion(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"subcommand", []string{"version"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset args before each test
			os.Args = append([]string{"pdfx"}, tt.args...)
			if err := runCLI(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestVersionFlag(t *testing.T) {
	if projectRoot == "" {
		t.Skip("project root not set")
	}

	for _, args := range [][]string{{"-v"}, {"--version"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			cmd := exec.Command("go", "run", filepath.Join(projectRoot, "main.go"))
			cmd.Args = append(cmd.Args, args...)
			cmd.Dir = projectRoot
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("command failed: %v\n%s", err, out)
			}
			if !strings.Contains(string(out), "v0.1.1") {
				t.Fatalf("expected version output, got: %s", out)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	if err := runCLI("merge", "-o", "merged.pdf"); err != nil {
		t.Fatalf("merge failed: %v", err)
	}
	if _, err := os.Stat("merged.pdf"); os.IsNotExist(err) {
		t.Fatal("merged.pdf was not created")
	}
}

func TestSplit(t *testing.T) {
	if err := runCLI("split", "-n", "test_file.pdf", "-f", "2"); err != nil {
		t.Fatalf("split failed: %v", err)
	}
	// pdfcpu naming convention for 2-page file: test_file_1.pdf, test_file_2.pdf
	if _, err := os.Stat("test_file_1.pdf"); os.IsNotExist(err) {
		t.Fatal("test_file_1.pdf was not created")
	}
	if _, err := os.Stat("test_file_2.pdf"); os.IsNotExist(err) {
		t.Fatal("test_file_2.pdf was not created")
	}
}

func TestExtract(t *testing.T) {
	if err := runCLI("extract", "-n", "test_file.pdf", "-p", "1", "-o", "extracted.pdf"); err != nil {
		t.Fatalf("extract failed: %v", err)
	}
	if _, err := os.Stat("extracted.pdf"); os.IsNotExist(err) {
		t.Fatal("extracted.pdf was not created")
	}
}

func TestMergeMissingDir(t *testing.T) {
	err := runCLI("merge", "-d", "./nonexistent", "-o", "out.pdf")
	if err == nil {
		t.Fatal("expected error for missing directory")
	}
}

func TestSplitMissingFile(t *testing.T) {
	err := runCLI("split", "-n", "nonexistent.pdf", "-f", "1")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestExtractMissingFile(t *testing.T) {
	err := runCLI("extract", "-n", "nonexistent.pdf", "-p", "1")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestExtractInvalidPage(t *testing.T) {
	err := runCLI("extract", "-n", "test_file.pdf", "-p", "abc")
	if err == nil {
		t.Fatal("expected error for invalid page")
	}
}
