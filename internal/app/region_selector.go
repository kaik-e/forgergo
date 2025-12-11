package app

import (
	"forger-companion/internal/config"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
)

type RegionSelector struct {
	app      *App
	window   fyne.Window
	callback func(*config.Region)
	
	startX, startY int
	endX, endY     int
	selecting      bool
}

func NewRegionSelector(app *App, callback func(*config.Region)) *RegionSelector {
	return &RegionSelector{
		app:      app,
		callback: callback,
	}
}

func (rs *RegionSelector) Show() {
	dialog := widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Region Selection"),
			widget.NewLabel("Instructions:"),
			widget.NewLabel("1. Click 'Start Selection'"),
			widget.NewLabel("2. Click and drag to select region"),
			widget.NewLabel("3. Release to confirm"),
			widget.NewSeparator(),
			widget.NewButton("Start Selection", rs.startSelection),
			widget.NewButton("Cancel", func() {
				// Close dialog
			}),
		),
		rs.app.window.Canvas(),
	)
	
	dialog.Show()
}

func (rs *RegionSelector) startSelection() {
	// Hide main window temporarily
	rs.app.window.Hide()
	
	// Simple region selection using mouse position
	// User needs to: click start position, then click end position
	
	// Get start position
	x, y := robotgo.Location()
	rs.startX, rs.startY = x, y
	
	// Wait for second click (simplified - in production use proper event handling)
	// For now, just create a simple dialog
	rs.app.window.Show()
	
	dialog := widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Click to set region corners"),
			widget.NewLabel("Top-Left: Click anywhere"),
			widget.NewButton("Set Top-Left", func() {
				x, y := robotgo.Location()
				rs.startX, rs.startY = x, y
			}),
			widget.NewLabel("Bottom-Right: Click anywhere"),
			widget.NewButton("Set Bottom-Right", func() {
				x, y := robotgo.Location()
				rs.endX, rs.endY = x, y
			}),
			widget.NewSeparator(),
			widget.NewButton("Confirm", func() {
				rs.confirmSelection()
			}),
			widget.NewButton("Cancel", func() {
				rs.app.window.Show()
			}),
		),
		rs.app.window.Canvas(),
	)
	
	dialog.Show()
}

func (rs *RegionSelector) confirmSelection() {
	width := rs.endX - rs.startX
	height := rs.endY - rs.startY
	
	if width < 0 {
		rs.startX, rs.endX = rs.endX, rs.startX
		width = -width
	}
	if height < 0 {
		rs.startY, rs.endY = rs.endY, rs.startY
		height = -height
	}
	
	region := &config.Region{
		X:      rs.startX,
		Y:      rs.startY,
		Width:  width,
		Height: height,
	}
	
	if rs.callback != nil {
		rs.callback(region)
	}
	
	rs.app.window.Show()
}

func (rs *RegionSelector) drawOverlay() *canvas.Rectangle {
	rect := canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 0, A: 128})
	rect.StrokeColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	rect.StrokeWidth = 2
	return rect
}
