#!/bin/bash
echo "=========================================="
echo "  FORGER COMPANION HYBRID BUILD"
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

# Copy Python companion app
echo "[*] Copying Python companion app..."
if [ ! -f "../forger-companion/main.py" ]; then
    echo "[ERROR] Python companion app not found!"
    exit 1
fi

# Create companion_app.py wrapper
cat > companion_app.py << 'EOF'
#!/usr/bin/env python3
"""Wrapper to run the companion app from the launcher directory"""
import sys
import os

# Add parent directory to path
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0, os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', 'forger-companion'))

# Import and run
from forger_companion.main import main
if __name__ == "__main__":
    main()
EOF

echo "[*] Created companion_app.py wrapper"
echo ""

echo "=========================================="
echo "  BUILD SUCCESSFUL!"
echo "=========================================="
echo ""
echo "  Launcher: ForgerCompanion.exe ($SIZE)"
echo "  Companion: companion_app.py"
echo ""
echo "  Distribution:"
echo "  1. Copy ForgerCompanion.exe"
echo "  2. Copy companion_app.py"
echo "  3. Copy ../forger-companion folder"
echo "  4. Zip and distribute!"
echo ""
echo "=========================================="
