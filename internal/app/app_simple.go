package app

import (
	"fmt"
	"forger-companion/internal/config"
	"forger-companion/internal/macro"
	"forger-companion/internal/ocr"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SimpleApp struct {
	cfg    *config.Config
	scanner *ocr.Scanner
	macro  *macro.Macro
	window fyne.Window
	
	statusLabel *widget.Label
	macroButton *widget.Button
	scanButton  *widget.Button
}

func NewSimple(cfg *config.Config) *SimpleApp {
	scanner := ocr.NewScanner()
	return &SimpleApp{
		cfg:     cfg,
		scanner: scanner,
		macro:   macro.New(cfg, scanner),
	}
}

func (a *SimpleApp) Run() {
	fyneApp := app.New()
	a.window = fyneApp.NewWindow("Forger Companion")
	a.window.Resize(fyne.NewSize(500, 400))
	
	// Title
	title := widget.NewLabelWithStyle(
		"🔨 Forger Companion",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	
	// Status
	a.statusLabel = widget.NewLabel("Ready")
	a.statusLabel.Wrapping = true
	
	// Buttons
	a.scanButton = widget.NewButton("Start Scan", a.toggleScan)
	a.macroButton = widget.NewButton("Start Macro", a.toggleMacro)
	settingsButton := widget.NewButton("Settings", a.openSettings)
	
	// Info
	infoLabel := widget.NewLabel(
		"Macro: Hold M1 at break position\n" +
		"Scan: Detect ores in selected region\n" +
		"Webhook: Send progress updates",
	)
	infoLabel.Wrapping = true
	
	// Layout
	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		infoLabel,
		widget.NewSeparator(),
		a.statusLabel,
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			a.scanButton,
			a.macroButton,
		),
		settingsButton,
	)
	
	a.window.SetContent(content)
	a.window.ShowAndRun()
	
	a.scanner.Close()
}

func (a *SimpleApp) toggleScan() {
	a.statusLabel.SetText("Scan not yet configured")
}

func (a *SimpleApp) toggleMacro() {
	if a.macro.IsRunning() {
		a.macro.Stop()
		a.macroButton.SetText("Start Macro")
		a.statusLabel.SetText("Macro stopped")
	} else {
		if err := a.macro.Start(); err != nil {
			a.statusLabel.SetText(fmt.Sprintf("Error: %v", err))
			return
		}
		a.macroButton.SetText("Stop Macro")
		a.statusLabel.SetText("Macro running...")
		
		// Update status periodically
		go func() {
			for a.macro.IsRunning() {
				time.Sleep(1 * time.Second)
			}
		}()
	}
}

func (a *SimpleApp) openSettings() {
	a.statusLabel.SetText("Settings: Edit ~/.forger-companion/settings.json")
}
