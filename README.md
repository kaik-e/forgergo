# Forger Companion (Go)

Lightweight Go rewrite of Forger Companion - **~15-30MB** vs 311MB Python version.

## Features

- ✅ OCR ore detection using Tesseract
- ✅ Real-time multiplier calculation
- ✅ Auto-scanning with configurable intervals
- ✅ Forge UI detection
- ✅ Macro functionality (auto-mining & selling)
- ✅ Webhook progress updates (Discord webhooks & bot DMs)
- ✅ Stats tracking (legendary ores, level, money)
- ✅ Region selection
- ✅ Tabbed UI (Calculator & Macro)

## Requirements

### Windows
- [Tesseract OCR](https://github.com/UB-Mannheim/tesseract/wiki) installed
  - Download: https://github.com/UB-Mannheim/tesseract/wiki
  - Add to PATH or install to default location

### macOS
```bash
brew install tesseract
```

### Linux
```bash
sudo apt-get install tesseract-ocr
```

## Building

```bash
# Install dependencies
go mod download

# Build
go build -o ForgerCompanion.exe .

# Build with optimizations (smaller binary)
go build -ldflags="-s -w" -o ForgerCompanion.exe .

# Further compress with UPX (optional)
upx --best --lzma ForgerCompanion.exe
```

## Development

```bash
# Run without building
go run .

# Run tests
go test ./...

# Format code
go fmt ./...
```

## Size Comparison

- **Python version**: 311 MB (with onnxruntime)
- **Go version**: ~15-30 MB (with Tesseract external)
- **Go version (UPX)**: ~8-15 MB

## Architecture

```
forger-companion-go/
├── main.go                 # Entry point
├── internal/
│   ├── app/               # GUI application
│   ├── config/            # Configuration management
│   ├── ocr/               # OCR scanning
│   ├── calculator/        # Forge calculations
│   ├── macro/             # Macro automation (TODO)
│   ├── webhook/           # Progress webhooks (TODO)
│   └── data/              # Ore data
└── go.mod
```

## Configuration

Config stored in `~/.forger-companion/settings.json`

Example:
```json
{
  "regions": {
    "ores_panel": {"x": 100, "y": 100, "width": 400, "height": 600}
  },
  "macro_buttons": {
    "break_position": {"x": 500, "y": 500},
    "inventory": {"key": "e"},
    "sell_tab": {"x": 300, "y": 200}
  },
  "webhook": {
    "enabled": true,
    "mode": "webhook",
    "webhook_url": "https://discord.com/api/webhooks/...",
    "cycle_interval": 5,
    "track_stats": true
  }
}
```

## TODO

- [ ] Hotkey support (F6 to toggle macro)
- [ ] System tray icon
- [ ] Auto-updater
- [ ] Better region selection UI
