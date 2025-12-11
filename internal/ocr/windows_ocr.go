// +build windows

package ocr

import (
	"forger-companion/internal/config"
	"forger-companion/internal/data"
	"image"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/kbinani/screenshot"
)

// WindowsScanner uses Windows built-in OCR (no CGO required)
type WindowsScanner struct {
}

func NewScanner() *Scanner {
	return &Scanner{
		client: &WindowsScanner{},
	}
}

type DetectedOre struct {
	Name       string
	Count      int
	Rarity     string
	Multiplier float64
}

type Scanner struct {
	client *WindowsScanner
}

func (s *Scanner) Close() {
	// Nothing to close
}

func (s *Scanner) CaptureRegion(region *config.Region) (image.Image, error) {
	bounds := image.Rect(region.X, region.Y, region.X+region.Width, region.Y+region.Height)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Use PowerShell OCR (Windows 10+)
func (s *Scanner) runOCR(img image.Image) (string, error) {
	// Save temp image
	// Run PowerShell OCR script
	// For now, return empty - this needs proper implementation
	log.Println("[OCR] Windows OCR not yet implemented - using pattern matching")
	return "", nil
}

func (s *Scanner) ScanForOres(region *config.Region) (map[string]DetectedOre, error) {
	img, err := s.CaptureRegion(region)
	if err != nil {
		return nil, err
	}

	text, _ := s.runOCR(img)
	return s.parseOres(text), nil
}

func (s *Scanner) parseOres(text string) map[string]DetectedOre {
	detected := make(map[string]DetectedOre)
	lines := strings.Split(text, "\n")

	countPattern := regexp.MustCompile(`x\s*(\d+)`)

	for i, line := range lines {
		lineLower := strings.ToLower(strings.TrimSpace(line))
		
		for oreName, oreData := range data.Ores {
			oreNameLower := strings.ToLower(oreName)
			baseName := strings.TrimSuffix(oreNameLower, " ore")
			
			if strings.Contains(lineLower, oreNameLower) || strings.Contains(lineLower, baseName) {
				count := 1
				
				for j := i; j < len(lines) && j < i+3; j++ {
					if matches := countPattern.FindStringSubmatch(lines[j]); len(matches) > 1 {
						if c, err := strconv.Atoi(matches[1]); err == nil && c > 0 && c < 100 {
							count = c
							break
						}
					}
				}
				
				detected[oreName] = DetectedOre{
					Name:       oreName,
					Count:      count,
					Rarity:     oreData.Rarity,
					Multiplier: oreData.Multiplier,
				}
				break
			}
		}
	}

	return detected
}

func (s *Scanner) DetectForgeUI(region *config.Region) (bool, bool, error) {
	img, err := s.CaptureRegion(region)
	if err != nil {
		return false, false, err
	}

	text, _ := s.runOCR(img)
	textLower := strings.ToLower(text)
	
	indicators := []string{"forge chances", "select ores", "forge!", "empty", "multiplier"}
	matches := 0
	for _, indicator := range indicators {
		if strings.Contains(textLower, indicator) {
			matches++
		}
	}

	isForgeUI := matches >= 2 || strings.Contains(textLower, "empty")
	emptyCount := strings.Count(textLower, "empty")
	hasOres := emptyCount < 4

	log.Printf("[OCR] ForgeUI=%v Empty=%d HasOres=%v", isForgeUI, emptyCount, hasOres)
	return isForgeUI, hasOres, nil
}

type Stats struct {
	LegendaryOres map[string]int
	Level         int
	Money         int
}

func (s *Scanner) ScanForStats(region *config.Region) (*Stats, error) {
	if region == nil {
		// Use full screen
		bounds := screenshot.GetDisplayBounds(0)
		region = &config.Region{
			X:      bounds.Min.X,
			Y:      bounds.Min.Y,
			Width:  bounds.Dx(),
			Height: bounds.Dy(),
		}
	}

	img, err := s.CaptureRegion(region)
	if err != nil {
		return nil, err
	}

	text, _ := s.runOCR(img)
	
	stats := &Stats{
		LegendaryOres: make(map[string]int),
	}

	textLower := strings.ToLower(text)
	lines := strings.Split(text, "\n")

	countPattern := regexp.MustCompile(`x\s*(\d+)`)
	for _, oreName := range data.LegendaryMythic {
		oreNameLower := strings.ToLower(oreName)
		if strings.Contains(textLower, oreNameLower) {
			for _, line := range lines {
				if strings.Contains(strings.ToLower(line), oreNameLower) {
					if matches := countPattern.FindStringSubmatch(line); len(matches) > 1 {
						if count, err := strconv.Atoi(matches[1]); err == nil {
							stats.LegendaryOres[oreName] += count
						}
					}
				}
			}
		}
	}

	levelPattern := regexp.MustCompile(`level\s+(\d+)`)
	if matches := levelPattern.FindStringSubmatch(textLower); len(matches) > 1 {
		if level, err := strconv.Atoi(matches[1]); err == nil {
			stats.Level = level
		}
	}

	moneyPattern := regexp.MustCompile(`\$\s*([\d,]+)`)
	if matches := moneyPattern.FindStringSubmatch(text); len(matches) > 1 {
		moneyStr := strings.ReplaceAll(matches[1], ",", "")
		if money, err := strconv.Atoi(moneyStr); err == nil {
			stats.Money = money
		}
	}

	return stats, nil
}

func runPowerShellOCR(imagePath string) (string, error) {
	script := `
Add-Type -AssemblyName System.Drawing
$img = [System.Drawing.Image]::FromFile($args[0])
# Windows OCR API would go here
$img.Dispose()
`
	cmd := exec.Command("powershell", "-Command", script, imagePath)
	output, err := cmd.Output()
	return string(output), err
}
