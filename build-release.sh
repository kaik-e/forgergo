#!/bin/bash
echo "=========================================="
echo "  FORGER COMPANION - RELEASE BUILD"
echo "=========================================="
echo ""

# Build launcher
echo "[*] Building launcher..."
go build -ldflags="-s -w -H windowsgui" -o ForgerCompanion.exe launcher.go

if [ $? -ne 0 ]; then
    echo "[ERROR] Launcher build failed!"
    exit 1
fi

SIZE=$(du -h ForgerCompanion.exe | cut -f1)
echo "[*] Launcher built: $SIZE"
echo ""

# Create release directory
RELEASE_DIR="ForgerCompanion-Release"
rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

echo "[*] Copying files..."

# Copy launcher
cp ForgerCompanion.exe "$RELEASE_DIR/"

# Copy Python companion app
if [ -d "../forger-companion" ]; then
    echo "[*] Copying Python companion app..."
    cp -r ../forger-companion "$RELEASE_DIR/companion-app"
    
    # Remove unnecessary files
    rm -rf "$RELEASE_DIR/companion-app/.git"
    rm -rf "$RELEASE_DIR/companion-app/__pycache__"
    rm -rf "$RELEASE_DIR/companion-app/.venv"
    rm -f "$RELEASE_DIR/companion-app/build.spec"
    rm -f "$RELEASE_DIR/companion-app/ForgerCompanion.exe"
else
    echo "[WARNING] Python companion app not found at ../forger-companion"
fi

# Create README
cat > "$RELEASE_DIR/README.txt" << 'EOF'
FORGER COMPANION
================

Installation:
1. Extract this folder
2. Double-click ForgerCompanion.exe

Requirements:
- Python 3.8+ (download from python.org if not installed)
- Windows 10/11

Features:
- Macro automation (auto-mining & selling)
- Webhook progress updates
- OCR ore detection
- Stats tracking

Configuration:
Settings saved in: %APPDATA%\.forger-companion\settings.json

Support:
Check the companion-app folder for source code
EOF

echo "[*] Created README.txt"
echo ""

echo "=========================================="
echo "  RELEASE READY!"
echo "=========================================="
echo ""
echo "  Directory: $RELEASE_DIR/"
echo "  Contents:"
echo "    - ForgerCompanion.exe (launcher)"
echo "    - companion-app/ (Python app)"
echo "    - README.txt"
echo ""
echo "  To distribute:"
echo "  1. Zip the folder"
echo "  2. Share with users"
echo ""
echo "=========================================="
