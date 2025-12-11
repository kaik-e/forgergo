package ocr

import (
	"forger-companion/internal/config"
	"forger-companion/internal/data"
	"image"
	"image/png"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/kbinani/screenshot"
	"github.com/otiai10/gosseract/v2"
)

type DetectedOre struct {
	Name       string
	Count      int
	Rarity     string
	Multiplier float64
}

type Scanner struct {
	client *gosseract.Client
}

func NewScanner() *Scanner {
	client := gosseract.NewClient()
	client.SetLanguage("eng")
	client.SetPageSegMode(gosseract.PSM_AUTO)
	
	return &Scanner{
		client: client,
	}
}

func (s *Scanner) Close() {
	if s.client != nil {
		s.client.Close()
	}
}

func (s *Scanner) CaptureRegion(region *config.Region) (image.Image, error) {
	bounds := image.Rect(region.X, region.Y, region.X+region.Width, region.Y+region.Height)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (s *Scanner) ScanForOres(region *config.Region) (map[string]DetectedOre, error) {
	img, err := s.CaptureRegion(region)
	if err != nil {
		return nil, err
	}

	// Save temp image for gosseract
	tmpFile := "temp_scan.png"
	if err := saveImage(img, tmpFile); err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile)

	s.client.SetImage(tmpFile)
	text, err := s.client.Text()
	if err != nil {
		return nil, err
	}

	return s.parseOres(text), nil
}

func (s *Scanner) parseOres(text string) map[string]DetectedOre {
	detected := make(map[string]DetectedOre)
	lines := strings.Split(text, "\n")

	// Pattern for ore counts (x1, x2, etc.)
	countPattern := regexp.MustCompile(`x\s*(\d+)`)

	for i, line := range lines {
		lineLower := strings.ToLower(strings.TrimSpace(line))
		
		// Check each ore name
		for oreName, oreData := range data.Ores {
			oreNameLower := strings.ToLower(oreName)
			baseName := strings.TrimSuffix(oreNameLower, " ore")
			
			if strings.Contains(lineLower, oreNameLower) || strings.Contains(lineLower, baseName) {
				count := 1
				
				// Look for count in nearby lines
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

	tmpFile := "temp_detect.png"
	if err := saveImage(img, tmpFile); err != nil {
		return false, false, err
	}
	defer os.Remove(tmpFile)

	s.client.SetImage(tmpFile)
	text, err := s.client.Text()
	if err != nil {
		return false, false, err
	}

	textLower := strings.ToLower(text)
	
	// Check for forge UI indicators
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
	img, err := s.CaptureRegion(region)
	if err != nil {
		return nil, err
	}

	tmpFile := "temp_stats.png"
	if err := saveImage(img, tmpFile); err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile)

	s.client.SetImage(tmpFile)
	text, err := s.client.Text()
	if err != nil {
		return nil, err
	}

	stats := &Stats{
		LegendaryOres: make(map[string]int),
	}

	textLower := strings.ToLower(text)
	lines := strings.Split(text, "\n")

	// Scan for legendary/mythic ores
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

	// Scan for level
	levelPattern := regexp.MustCompile(`level\s+(\d+)`)
	if matches := levelPattern.FindStringSubmatch(textLower); len(matches) > 1 {
		if level, err := strconv.Atoi(matches[1]); err == nil {
			stats.Level = level
		}
	}

	// Scan for money
	moneyPattern := regexp.MustCompile(`\$\s*([\d,]+)`)
	if matches := moneyPattern.FindStringSubmatch(text); len(matches) > 1 {
		moneyStr := strings.ReplaceAll(matches[1], ",", "")
		if money, err := strconv.Atoi(moneyStr); err == nil {
			stats.Money = money
		}
	}

	return stats, nil
}

func saveImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	
	return png.Encode(f, img)
}
