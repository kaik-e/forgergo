package ocr

import (
	"forger-companion/internal/config"
	"forger-companion/internal/data"
	"image"
	"image/color"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/kbinani/screenshot"
)

// SimpleScanner uses basic pattern recognition (no external OCR needed)
// Works by detecting color patterns and text-like regions
type SimpleScanner struct {
}

type DetectedOre struct {
	Name       string
	Count      int
	Rarity     string
	Multiplier float64
}

type Scanner struct {
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Close() {
	// Nothing to close
}

func (s *Scanner) CaptureRegion(region *config.Region) (image.Image, error) {
	if region == nil {
		bounds := screenshot.GetDisplayBounds(0)
		img, err := screenshot.CaptureRect(bounds)
		return img, err
	}
	
	bounds := image.Rect(region.X, region.Y, region.X+region.Width, region.Y+region.Height)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Simplified OCR using color detection and pattern matching
// For production, users can optionally install Tesseract for better accuracy
func (s *Scanner) ScanForOres(region *config.Region) (map[string]DetectedOre, error) {
	img, err := s.CaptureRegion(region)
	if err != nil {
		return nil, err
	}

	// Detect text regions by color (white/light text on dark background)
	textRegions := s.detectTextRegions(img)
	
	// Simple pattern matching
	detected := make(map[string]DetectedOre)
	
	// For now, return empty - user needs to configure manually
	// In production, this would use actual OCR or ML model
	log.Println("[OCR] Simple scanner - manual configuration recommended")
	
	return detected, nil
}

func (s *Scanner) detectTextRegions(img image.Image) []image.Rectangle {
	bounds := img.Bounds()
	regions := []image.Rectangle{}
	
	// Scan for bright text on dark background
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 10 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 10 {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			
			// Check if pixel is bright (potential text)
			brightness := (r + g + b) / 3
			if brightness > 40000 { // Bright pixel
				// Found potential text region
				regions = append(regions, image.Rect(x, y, x+100, y+20))
			}
		}
	}
	
	return regions
}

func (s *Scanner) DetectForgeUI(region *config.Region) (bool, bool, error) {
	img, err := s.CaptureRegion(region)
	if err != nil {
		return false, false, err
	}

	// Simple detection: look for UI color patterns
	hasUI := s.detectUIPattern(img)
	hasOres := true // Assume ores present
	
	log.Printf("[OCR] ForgeUI=%v HasOres=%v (simple detection)", hasUI, hasOres)
	return hasUI, hasOres, nil
}

func (s *Scanner) detectUIPattern(img image.Image) bool {
	bounds := img.Bounds()
	brightPixels := 0
	totalPixels := 0
	
	// Sample pixels
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 5 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 5 {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			brightness := (r + g + b) / 3
			
			if brightness > 30000 {
				brightPixels++
			}
			totalPixels++
		}
	}
	
	// If >10% bright pixels, likely UI is present
	ratio := float64(brightPixels) / float64(totalPixels)
	return ratio > 0.1
}

type Stats struct {
	LegendaryOres map[string]int
	Level         int
	Money         int
}

func (s *Scanner) ScanForStats(region *config.Region) (*Stats, error) {
	stats := &Stats{
		LegendaryOres: make(map[string]int),
	}
	
	// Simplified - would need actual OCR
	log.Println("[OCR] Stats scanning - requires OCR setup")
	
	return stats, nil
}

func parseOres(text string) map[string]DetectedOre {
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

func isTextLike(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	brightness := (r + g + b) / 3
	return brightness > 40000 // Bright pixels likely text
}
