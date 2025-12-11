package app

import (
	"fmt"
	"forger-companion/internal/calculator"
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

type App struct {
	cfg     *config.Config
	scanner *ocr.Scanner
	macro   *macro.Macro
	window  fyne.Window
	
	// UI elements
	multiplierLabel *widget.Label
	oresLabel       *widget.Label
	statusLabel     *widget.Label
	scanButton      *widget.Button
	macroButton     *widget.Button
	
	// State
	scanning bool
	stopChan chan bool
}

func New(cfg *config.Config) *App {
	scanner := ocr.NewScanner()
	return &App{
		cfg:      cfg,
		scanner:  scanner,
		macro:    macro.New(cfg, scanner),
		stopChan: make(chan bool),
	}
}

func (a *App) Run() {
	fyneApp := app.New()
	a.window = fyneApp.NewWindow("Forger Companion")
	
	a.buildUI()
	
	// Set window properties
	a.window.Resize(fyne.NewSize(400, 300))
	if alwaysOnTop, ok := a.cfg.Preferences["always_on_top"].(bool); ok && alwaysOnTop {
		a.window.SetOnTop(true)
	}
	
	a.window.ShowAndRun()
	
	// Cleanup
	a.scanner.Close()
}

func (a *App) buildUI() {
	// Title
	title := widget.NewLabelWithStyle("Forger Companion", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	
	// Multiplier display
	a.multiplierLabel = widget.NewLabelWithStyle("Multiplier: 1.00x", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	
	// Ores display
	a.oresLabel = widget.NewLabel("No ores detected")
	a.oresLabel.Wrapping = fyne.TextWrapWord
	
	// Status
	a.statusLabel = widget.NewLabel("Ready")
	
	// Buttons
	regionButton := widget.NewButton("Select Region", a.selectRegion)
	a.scanButton = widget.NewButton("Start Scan", a.toggleScan)
	a.macroButton = widget.NewButton("Start Macro", a.toggleMacro)
	
	// Tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Calculator", container.NewVBox(
			a.multiplierLabel,
			widget.NewSeparator(),
			a.oresLabel,
			widget.NewSeparator(),
			container.NewGridWithColumns(2,
				regionButton,
				a.scanButton,
			),
		)),
		container.NewTabItem("Macro", container.NewVBox(
			widget.NewLabel("Macro Settings"),
			widget.NewLabel("Configure macro buttons in settings"),
			widget.NewSeparator(),
			a.macroButton,
		)),
	)
	
	// Layout
	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		a.statusLabel,
		widget.NewSeparator(),
		tabs,
	)
	
	a.window.SetContent(content)
}

func (a *App) selectRegion() {
	selector := NewRegionSelector(a, func(region *config.Region) {
		a.cfg.Regions["ores_panel"] = region
		a.cfg.Save()
		a.statusLabel.SetText("Region saved!")
	})
	selector.Show()
}

func (a *App) toggleScan() {
	if a.scanning {
		a.stopScan()
	} else {
		a.startScan()
	}
}

func (a *App) toggleMacro() {
	if a.macro.IsRunning() {
		a.macro.Stop()
		a.macroButton.SetText("Start Macro")
		a.statusLabel.SetText("Macro stopped")
	} else {
		if err := a.macro.Start(); err != nil {
			a.statusLabel.SetText(fmt.Sprintf("Macro error: %v", err))
			return
		}
		a.macroButton.SetText("Stop Macro")
		a.statusLabel.SetText("Macro running...")
	}
}

func (a *App) startScan() {
	region := a.cfg.Regions["ores_panel"]
	if region == nil {
		a.statusLabel.SetText("Please select a region first")
		return
	}
	
	a.scanning = true
	a.scanButton.SetText("Stop Scan")
	a.statusLabel.SetText("Scanning...")
	
	go a.scanLoop(region)
}

func (a *App) stopScan() {
	a.scanning = false
	a.scanButton.SetText("Start Scan")
	a.statusLabel.SetText("Stopped")
	a.stopChan <- true
}

func (a *App) scanLoop(region *config.Region) {
	interval := 2 * time.Second
	if scanInterval, ok := a.cfg.Preferences["scan_interval"].(float64); ok {
		interval = time.Duration(scanInterval * float64(time.Second))
	}
	
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	for a.scanning {
		select {
		case <-ticker.C:
			a.performScan(region)
		case <-a.stopChan:
			return
		}
	}
}

func (a *App) performScan(region *config.Region) {
	// Check if forge UI is open
	isForgeUI, hasOres, err := a.scanner.DetectForgeUI(region)
	if err != nil {
		log.Printf("Error detecting forge UI: %v", err)
		return
	}
	
	if !isForgeUI || !hasOres {
		a.multiplierLabel.SetText("Multiplier: 1.00x")
		a.oresLabel.SetText("Forge UI not detected or no ores placed")
		return
	}
	
	// Scan for ores
	ores, err := a.scanner.ScanForOres(region)
	if err != nil {
		log.Printf("Error scanning ores: %v", err)
		return
	}
	
	if len(ores) == 0 {
		a.oresLabel.SetText("No ores detected")
		return
	}
	
	// Calculate multiplier
	result := calculator.Calculate(ores)
	
	// Update UI
	a.multiplierLabel.SetText(fmt.Sprintf("Multiplier: %.2fx", result.TotalMultiplier))
	
	oresText := fmt.Sprintf("Detected %d ores:\n", len(ores))
	for _, ore := range ores {
		oresText += fmt.Sprintf("• %s x%d (%.1fx)\n", ore.Name, ore.Count, ore.Multiplier)
	}
	a.oresLabel.SetText(oresText)
	
	a.statusLabel.SetText(fmt.Sprintf("Last scan: %s", time.Now().Format("15:04:05")))
}
