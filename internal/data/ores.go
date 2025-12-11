package data

type Ore struct {
	Name       string
	Rarity     string
	Multiplier float64
}

var Ores = map[string]Ore{
	"Coal Ore":         {"Coal Ore", "common", 1.0},
	"Copper Ore":       {"Copper Ore", "common", 1.1},
	"Iron Ore":         {"Iron Ore", "common", 1.2},
	"Tin Ore":          {"Tin Ore", "uncommon", 1.3},
	"Silver Ore":       {"Silver Ore", "uncommon", 1.4},
	"Gold Ore":         {"Gold Ore", "uncommon", 1.5},
	"Topaz Ore":        {"Topaz Ore", "rare", 1.6},
	"Emerald Ore":      {"Emerald Ore", "rare", 1.7},
	"Ruby Ore":         {"Ruby Ore", "rare", 1.8},
	"Sapphire Ore":     {"Sapphire Ore", "legendary", 2.0},
	"Titanium Ore":     {"Titanium Ore", "legendary", 2.2},
	"Orichalcum Ore":   {"Orichalcum Ore", "legendary", 2.4},
	"Mythril Ore":      {"Mythril Ore", "mythical", 2.6},
	"Adamantite Ore":   {"Adamantite Ore", "mythical", 2.8},
	"Eye Ore":          {"Eye Ore", "epic", 1.9},
	"Rivalite Ore":     {"Rivalite Ore", "rare", 1.75},
	"Magmaite Ore":     {"Magmaite Ore", "epic", 1.95},
}

var LegendaryMythic = []string{
	"Sapphire Ore",
	"Titanium Ore",
	"Orichalcum Ore",
	"Mythril Ore",
	"Adamantite Ore",
}
