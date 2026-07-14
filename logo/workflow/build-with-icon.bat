@echo off
setlocal enabledelayedexpansion

cd /d "%~dp0..\.."

if not exist "logo\pdfx.ico" (
    echo [icon] logo/pdfx.ico missing, converting from JPG
    python -c "from PIL import Image; Image.open('logo/pdfx.jpg').save('logo/pdfx.ico', format='ICO', sizes=[(256,256)])"
) else (
    for %%A in ("logo\pdfx.jpg") do set "jpg_time=%%~tA"
    for %%A in ("logo\pdfx.ico") do set "ico_time=%%~tA"
    echo [icon] JPG: !jpg_time!
    echo [icon] ICO: !ico_time!
    echo [icon] logo/pdfx.ico exists, assuming up to date
)

echo [res] Generating pdfx.rc
(
echo 1 ICON "logo/pdfx.ico"
) > pdfx.rc

echo [res] Compiling pdfx.syso
windres -o pdfx.syso pdfx.rc
if errorlevel 1 (
    echo [error] windres failed
    exit /b 1
)

echo [build] Building pdfx.exe
go build -o pdfx.exe .
if errorlevel 1 (
    echo [error] go build failed
    exit /b 1
)

echo [done] pdfx.exe built successfully with icon
pause
