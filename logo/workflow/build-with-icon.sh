#!/usr/bin/env bash
set -euo pipefail

# Resolve project root (script lives in logo/workflow/)
PROJECT_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$PROJECT_ROOT"

# 1. Convert JPG to ICO when missing or outdated
if [ ! -f "logo/pdfx.ico" ] || [ "logo/pdfx.jpg" -nt "logo/pdfx.ico" ]; then
    echo "[icon] Converting logo/pdfx.jpg -> logo/pdfx.ico"
    python -c "from PIL import Image; Image.open('logo/pdfx.jpg').save('logo/pdfx.ico', format='ICO', sizes=[(256,256)])"
else
    echo "[icon] logo/pdfx.ico is up to date"
fi

# 2. Generate resource script
cat > pdfx.rc << 'EOF'
1 ICON "logo/pdfx.ico"
EOF
echo "[res] Generated pdfx.rc"

# 3. Compile Windows resource
windres -o pdfx.syso pdfx.rc
echo "[res] Compiled pdfx.syso"

# 4. Build exe
go build -o pdfx.exe .
echo "[build] pdfx.exe built successfully"
