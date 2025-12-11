package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Region struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type MacroButton struct {
	X   *int    `json:"x,omitempty"`
	Y   *int    `json:"y,omitempty"`
	Key *string `json:"key,omitempty"`
}

type WebhookSettings struct {
	Enabled       bool   `json:"enabled"`
	Mode          string `json:"mode"` // "bot" or "webhook"
	DiscordID     string `json:"discord_id"`
	WebhookURL    string `json:"webhook_url"`
	CycleInterval int    `json:"cycle_interval"`
	SendGIF       bool   `json:"send_gif"`
	GIFFrames     int    `json:"gif_frames"`
	GIFDuration   int    `json:"gif_duration"`
	TrackStats    bool   `json:"track_stats"`
}

type Config struct {
	SetupComplete bool                       `json:"setup_complete"`
	Regions       map[string]*Region         `json:"regions"`
	MacroButtons  map[string]*MacroButton    `json:"macro_buttons"`
	MacroSettings map[string]interface{}     `json:"macro_settings"`
	Webhook       WebhookSettings            `json:"webhook"`
	Preferences   map[string]interface{}     `json:"preferences"`
	Window        map[string]interface{}     `json:"window"`
}

func Default() *Config {
	return &Config{
		SetupComplete: false,
		Regions:       make(map[string]*Region),
		MacroButtons:  make(map[string]*MacroButton),
		MacroSettings: map[string]interface{}{
			"enabled":       false,
			"hold_duration": 5,
			"auto_sell":     true,
		},
		Webhook: WebhookSettings{
			Enabled:       false,
			Mode:          "bot",
			CycleInterval: 5,
			SendGIF:       false,
			GIFFrames:     5,
			GIFDuration:   500,
			TrackStats:    false,
		},
		Preferences: map[string]interface{}{
			"auto_mode":        true,
			"always_on_top":    true,
			"auto_switch_tab":  true,
			"opacity":          95,
			"scan_interval":    2.0,
			"macro_hotkey":     "f6",
		},
		Window: make(map[string]interface{}),
	}
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".forger-companion", "settings.json")
}

func Load() (*Config, error) {
	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	path := configPath()
	dir := filepath.Dir(path)
	
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
