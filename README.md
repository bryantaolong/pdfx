# pdfx

A Git-like CLI tool for merging, splitting, and extracting PDF pages.

## Build

```bash
go mod tidy
go build -o pdfx .
```

## Usage

```bash
./pdfx merge --dir/-d <directory> --output/-o <output.pdf>
./pdfx split --name/-n <input.pdf> --from/-f <page>
./pdfx extract --name/-n <input.pdf> --pages/-p <1,2,3> --output/-o <output.pdf>
```

### Commands

| Command | Description |
|---------|-------------|
| `merge` | Merge all PDF files in a directory into one file |
| `split` | Split a PDF into two files at a specified page number |
| `extract` | Extract specified pages from a PDF and merge them into a new file |

### Flags

- `--dir, -d`: Directory containing PDF files (default: current directory)
- `--output, -o`: Output file path (merge, extract)
- `--name, -n`: Input PDF file path (split, extract)
- `--from, -f`: Start page of the second file, 1-based (split)
- `--pages, -p`: Comma-separated page numbers, e.g. 1,2,3,4 (extract)
