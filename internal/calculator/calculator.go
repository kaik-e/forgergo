package calculator

import (
	"forger-companion/internal/ocr"
	"math"
)

type Result struct {
	TotalMultiplier float64
	OreCount        int
	Ores            map[string]ocr.DetectedOre
}

func Calculate(ores map[string]ocr.DetectedOre) *Result {
	totalMultiplier := 1.0
	totalOres := 0

	for _, ore := range ores {
		totalMultiplier *= math.Pow(ore.Multiplier, float64(ore.Count))
		totalOres += ore.Count
	}

	return &Result{
		TotalMultiplier: totalMultiplier,
		OreCount:        totalOres,
		Ores:            ores,
	}
}
