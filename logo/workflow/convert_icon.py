#!/usr/bin/env python3
"""Convert logo/pdfx.jpg to logo/pdfx.ico for Windows executable icon."""

import sys
from pathlib import Path

try:
    from PIL import Image
except ImportError:
    print("[error] Pillow is not installed. Run: pip install pillow")
    sys.exit(1)


def convert(jpg_path: str, ico_path: str, size: int = 256) -> None:
    jpg = Path(jpg_path)
    ico = Path(ico_path)

    if not jpg.exists():
        print(f"[error] Source image not found: {jpg}")
        sys.exit(1)

    # Only convert if JPG is newer or ICO does not exist
    if ico.exists() and ico.stat().st_mtime >= jpg.stat().st_mtime:
        print(f"[icon] {ico} is up to date")
        return

    print(f"[icon] Converting {jpg} -> {ico}")
    img = Image.open(jpg)
    img.save(ico, format="ICO", sizes=[(size, size)])
    print(f"[icon] Done: {ico}")


if __name__ == "__main__":
    project_root = Path(__file__).resolve().parent.parent.parent
    convert(
        jpg_path=str(project_root / "logo" / "pdfx.jpg"),
        ico_path=str(project_root / "logo" / "pdfx.ico"),
    )
