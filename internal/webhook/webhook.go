package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forger-companion/internal/config"
	"forger-companion/internal/ocr"
	"image"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/kbinani/screenshot"
)

type Manager struct {
	cfg *config.Config
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

func (m *Manager) ShouldSendUpdate(cycle int) bool {
	if !m.cfg.Webhook.Enabled {
		return false
	}
	if m.cfg.Webhook.Mode == "webhook" && m.cfg.Webhook.WebhookURL == "" {
		return false
	}
	if m.cfg.Webhook.Mode == "bot" && m.cfg.Webhook.DiscordID == "" {
		return false
	}
	return cycle > 0 && cycle%m.cfg.Webhook.CycleInterval == 0
}

func (m *Manager) TrackStats() bool {
	return m.cfg.Webhook.TrackStats
}

func (m *Manager) SendUpdate(cycle int, stats *ocr.Stats) error {
	if m.cfg.Webhook.Mode == "webhook" {
		return m.sendWebhook(cycle, stats)
	}
	return m.sendBotDM(cycle, stats)
}

func (m *Manager) sendWebhook(cycle int, stats *ocr.Stats) error {
	// Capture screenshot
	img, err := m.captureScreen()
	if err != nil {
		return err
	}

	// Create embed
	embed := map[string]interface{}{
		"title": "🔨 Macro Progress Update",
		"color": 5793522,
		"fields": []map[string]interface{}{
			{"name": "Cycle", "value": fmt.Sprintf("#%d", cycle), "inline": true},
			{"name": "Time", "value": fmt.Sprintf("<t:%d:R>", time.Now().Unix()), "inline": true},
		},
		"image": map[string]string{
			"url": fmt.Sprintf("attachment://progress_cycle_%d.png", cycle),
		},
		"footer": map[string]string{
			"text": "Forger Companion",
		},
	}

	// Add stats if provided
	if stats != nil {
		fields := embed["fields"].([]map[string]interface{})
		
		if len(stats.LegendaryOres) > 0 {
			oresText := ""
			for name, count := range stats.LegendaryOres {
				oresText += fmt.Sprintf("• %s: %d\n", name, count)
			}
			fields = append(fields, map[string]interface{}{
				"name":   "🌟 Legendary/Mythic Ores",
				"value":  oresText,
				"inline": false,
			})
		}
		
		if stats.Level > 0 {
			fields = append(fields, map[string]interface{}{
				"name":   "📊 Level",
				"value":  fmt.Sprintf("%d", stats.Level),
				"inline": true,
			})
		}
		
		if stats.Money > 0 {
			fields = append(fields, map[string]interface{}{
				"name":   "💰 Money",
				"value":  fmt.Sprintf("$%d", stats.Money),
				"inline": true,
			})
		}
		
		embed["fields"] = fields
	}

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add image
	imgBuf := &bytes.Buffer{}
	if err := png.Encode(imgBuf, img); err != nil {
		return err
	}
	
	part, err := writer.CreateFormFile("file", fmt.Sprintf("progress_cycle_%d.png", cycle))
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, imgBuf); err != nil {
		return err
	}

	// Add payload
	payload := map[string]interface{}{
		"embeds": []interface{}{embed},
	}
	payloadJSON, _ := json.Marshal(payload)
	writer.WriteField("payload_json", string(payloadJSON))
	writer.Close()

	// Send request
	req, err := http.NewRequest("POST", m.cfg.Webhook.WebhookURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("webhook failed: %d", resp.StatusCode)
	}

	log.Println("[Webhook] Update sent successfully!")
	return nil
}

func (m *Manager) sendBotDM(cycle int, stats *ocr.Stats) error {
	// Capture screenshot
	img, err := m.captureScreen()
	if err != nil {
		return err
	}

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add image
	imgBuf := &bytes.Buffer{}
	if err := png.Encode(imgBuf, img); err != nil {
		return err
	}
	
	part, err := writer.CreateFormFile("image", fmt.Sprintf("progress_cycle_%d.png", cycle))
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, imgBuf); err != nil {
		return err
	}

	// Add form fields
	writer.WriteField("discord_id", m.cfg.Webhook.DiscordID)
	writer.WriteField("cycle", fmt.Sprintf("%d", cycle))
	writer.WriteField("timestamp", time.Now().Format(time.RFC3339))
	
	// TODO: Add license_key from auth
	
	writer.Close()

	// Send to bot API
	botURL := "https://forger-production.up.railway.app/api/progress"
	req, err := http.NewRequest("POST", botURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("bot API failed: %d", resp.StatusCode)
	}

	log.Println("[Webhook] Bot DM sent successfully!")
	return nil
}

func (m *Manager) captureScreen() (image.Image, error) {
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}
	
	// Resize to half for smaller file size
	resized := image.NewRGBA(image.Rect(0, 0, bounds.Dx()/2, bounds.Dy()/2))
	// Simple nearest-neighbor resize
	for y := 0; y < bounds.Dy()/2; y++ {
		for x := 0; x < bounds.Dx()/2; x++ {
			resized.Set(x, y, img.At(x*2, y*2))
		}
	}
	
	return resized, nil
}
