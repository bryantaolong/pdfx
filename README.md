<div align="center">
  <a href="https://github.com/bryantaolong/pdfx">
    <img width="200" height="200" src="logo/pdfx.jpg" alt="pdfx logo"/>
  </a>
</div>

**PDFX** empowers you to merge, split, extract and remove PDF pages right in your terminal.

---

## ✨ Features

* Merge all PDF files in a directory into one file
* Split a PDF into two files at a specified page number
* Extract specified pages from a PDF and merge them into a new file
* Remove specified pages from a PDF and save the result to a new file
* Check version via `pdfx version` or `pdfx -v`
* Color-coded output for easy operation feedback

---

## 🚀 Installation & Build

### 1. Install Go

Make sure Go is installed and environment variables are set:

```bash
go version
```

### 2. Clone the repository

```bash
git clone https://github.com/bryantaolong/pdfx.git
cd pdfx
```

### 3. Build the executable

#### Windows

```powershell
go build -o pdfx.exe
```

#### Linux / macOS

```bash
go build -o pdfx
```

> Optional: Add the executable to your system PATH for global usage.

---

## 🎮 Usage

### Help

```bash
pdfx --help
pdfx merge --help
pdfx split --help
pdfx extract --help
pdfx remove --help
```

### Merge PDFs

Merge all PDF files in a directory into one file:

```bash
# Merge all PDFs in current directory
pdfx merge --output merged.pdf

# Merge PDFs from a specific directory
pdfx merge -d ./path/to/pdfs -o merged.pdf
```

### Split PDF

Split a PDF into two files at a specified page number:

```bash
pdfx split input.pdf --from 10
```

This splits `input.pdf` into:
* `input_1-9.pdf` (pages 1-9)
* `input_10-56.pdf` (pages 10 to end, actual end page depends on input)

### Extract PDF

Extract specified pages from a PDF and merge them into a new file:

```bash
pdfx extract input.pdf --pages 1,3,5 --output extracted.pdf
```

### Remove PDF Pages

Remove specified pages from a PDF and save the result to a new file:

```bash
pdfx remove input.pdf --pages 1,3,5 --output removed.pdf
```

---

## Example

```bash
# Check version
pdfx -v

# Merge all PDFs in a directory
pdfx merge -d ./docs -o combined.pdf

# Split a PDF at page 15
pdfx split book.pdf -f 15

# Extract pages 1, 2, and 3
pdfx extract book.pdf -p 1,2,3 -o chapters.pdf

# Remove pages 1, 2, and 3
pdfx remove book.pdf -p 1,2,3 -o remaining.pdf
```

---

## Project Structure

```
pdfx/
├─ cmd/           # Cobra command modules
│  ├─ root.go
│  ├─ version.go
│  └─ commands/
│     ├─ util.go
│     ├─ merge.go
│     ├─ split.go
│     ├─ extract.go
│     └─ remove.go
├─ test/          # Integration tests
├─ logo/          # Logo assets and icon workflow
├─ main.go
└─ README.md
```

---

## 💾 Notes

* Page numbers are 1-based in all commands
* The `--from` flag in `split` indicates the start page of the second file
* The `--pages` flag in `extract` accepts comma-separated page numbers
* The `--pages` flag in `remove` accepts comma-separated page numbers to remove
* Removing all pages is not allowed and will result in an error

---

## License

[MIT License](./LICENSE)
