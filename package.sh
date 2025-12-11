#!/bin/bash
echo "=========================================="
echo "  PACKAGING FORGER COMPANION"
echo "=========================================="
echo ""

if [ ! -f "ForgerCompanion.exe" ]; then
    echo "[ERROR] ForgerCompanion.exe not found!"
    echo "Run ./build-msys2.sh first"
    exit 1
fi

echo "[*] Creating package directory..."
mkdir -p ForgerCompanion-Package
cd ForgerCompanion-Package

echo "[*] Copying executable..."
cp ../ForgerCompanion.exe .

echo "[*] Copying Tesseract runtime..."
cp -r /mingw64/bin/libtesseract*.dll .
cp -r /mingw64/bin/leptonica*.dll .
cp -r /mingw64/bin/libpng*.dll .
cp -r /mingw64/bin/libjpeg*.dll .
cp -r /mingw64/bin/libwebp*.dll .
cp -r /mingw64/bin/zlib*.dll .
cp -r /mingw64/bin/libgcc*.dll .
cp -r /mingw64/bin/libstdc++*.dll .
cp -r /mingw64/bin/libwinpthread*.dll .

echo "[*] Copying Tesseract data..."
mkdir -p tessdata
cp -r /mingw64/share/tesseract/tessdata/* tessdata/ 2>/dev/null || true

echo "[*] Creating README..."
cat > README.txt << 'EOF'
FORGER COMPANION GO
===================

Installation:
1. Extract this folder anywhere
2. Double-click ForgerCompanion.exe

No additional installation required!

Features:
- OCR ore detection
- Real-time multiplier calculation
- Macro automation
- Webhook progress updates
- Stats tracking

Configuration:
Settings are saved in: %APPDATA%\.forger-companion\settings.json

Troubleshooting:
If you get a DLL error, make sure all files are in the same folder.
EOF

cd ..

echo ""
echo "=========================================="
echo "  PACKAGING COMPLETE!"
echo "=========================================="
echo ""
echo "  Package: ForgerCompanion-Package/"
echo "  Ready to distribute!"
echo ""
