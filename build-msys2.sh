#!/bin/bash
echo "=========================================="
echo "  FORGER COMPANION GO - BUILD (MSYS2)"
echo "=========================================="
echo ""

# Check if in MSYS2
if [ -z "$MSYSTEM" ]; then
    echo "[ERROR] This script must be run in MSYS2 MinGW 64-bit terminal"
    echo ""
    echo "Please:"
    echo "1. Open 'MSYS2 MinGW 64-bit' from Start Menu"
    echo "2. Navigate to this directory"
    echo "3. Run: ./build-msys2.sh"
    exit 1
fi

echo "[*] Checking dependencies..."

# Check if packages are installed
if ! command -v gcc &> /dev/null; then
    echo "[*] Installing gcc..."
    pacman -S --noconfirm mingw-w64-x86_64-gcc
fi

if ! command -v go &> /dev/null; then
    echo "[*] Installing Go..."
    pacman -S --noconfirm mingw-w64-x86_64-go
fi

if ! pkg-config --exists tesseract; then
    echo "[*] Installing Tesseract..."
    pacman -S --noconfirm mingw-w64-x86_64-tesseract-ocr mingw-w64-x86_64-leptonica
fi

if ! command -v pkg-config &> /dev/null; then
    echo "[*] Installing pkg-config..."
    pacman -S --noconfirm mingw-w64-x86_64-pkg-config
fi

echo ""
echo "[*] Dependencies installed!"
echo ""
echo "[*] Downloading Go modules..."
export PKG_CONFIG_PATH="/mingw64/lib/pkgconfig"
export CGO_ENABLED=1
go mod tidy

echo ""
echo "[*] Building ForgerCompanion.exe..."
echo "[*] This may take 2-5 minutes on first build..."
go build -v -ldflags="-s -w -H windowsgui" -o ForgerCompanion.exe .

if [ $? -eq 0 ]; then
    SIZE=$(du -h ForgerCompanion.exe | cut -f1)
    echo ""
    echo "=========================================="
    echo "  BUILD SUCCESSFUL!"
    echo "=========================================="
    echo ""
    echo "  Output: ForgerCompanion.exe"
    echo "  Size: $SIZE"
    echo ""
    echo "[*] Bundling required DLLs..."
    
    # Copy required DLLs (find them dynamically)
    cp /mingw64/bin/libtesseract*.dll . 2>/dev/null || true
    cp /mingw64/bin/leptonica*.dll . 2>/dev/null || true
    cp /mingw64/bin/libgomp*.dll . 2>/dev/null || true
    cp /mingw64/bin/libgcc*.dll . 2>/dev/null || true
    cp /mingw64/bin/libstdc++*.dll . 2>/dev/null || true
    cp /mingw64/bin/libwinpthread*.dll . 2>/dev/null || true
    cp /mingw64/bin/libpng*.dll . 2>/dev/null || true
    cp /mingw64/bin/libjpeg*.dll . 2>/dev/null || true
    cp /mingw64/bin/libwebp*.dll . 2>/dev/null || true
    cp /mingw64/bin/zlib*.dll . 2>/dev/null || true
    
    echo "[*] DLLs bundled!"
    ls -lh *.dll 2>/dev/null | wc -l | xargs echo "[*] Total DLLs:"
    echo ""
    echo "  Files ready to distribute:"
    echo "  - ForgerCompanion.exe"
    echo "  - *.dll files"
    echo ""
    echo "=========================================="
else
    echo ""
    echo "[ERROR] Build failed!"
    exit 1
fi
