# Icon Workflow

Automated steps to rebuild `pdfx.exe` with the current `logo/pdfx.jpg` icon.

## Prerequisites

- Git Bash / MSYS2 shell (provides `bash`, `windres`)  **or**  Windows Command Prompt
- Python with Pillow (`pip install pillow`)
- Go toolchain
- `windres` available in PATH (e.g. from MinGW)

## Scripts

- **`build-with-icon.sh`** — Git Bash / MSYS2 / WSL
- **`build-with-icon.bat`** — Windows Command Prompt / double-click

## Usage

### Bash (Git Bash / MSYS2)

```bash
cd logo/workflow
bash build-with-icon.sh
```

### Windows CMD / double-click

```
logo\workflow\build-with-icon.bat
```

## What it does

1. Converts `logo/pdfx.jpg` to `logo/pdfx.ico` when the JPG is newer (bash version only).
2. Writes `pdfx.rc` in the project root.
3. Compiles `pdfx.syso` with `windres`.
4. Builds `pdfx.exe` with the embedded icon.

## Generated files (project root)

- `logo/pdfx.ico`
- `pdfx.rc`
- `pdfx.syso`
- `pdfx.exe`
