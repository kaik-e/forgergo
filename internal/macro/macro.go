package macro

import (
	"forger-companion/internal/config"
	"forger-companion/internal/ocr"
	"forger-companion/internal/webhook"
	"log"
	"time"

	"github.com/go-vgo/robotgo"
)

type Macro struct {
	cfg            *config.Config
	running        bool
	stopChan       chan bool
	webhookManager *webhook.Manager
	scanner        *ocr.Scanner
}

func New(cfg *config.Config, scanner *ocr.Scanner) *Macro {
	return &Macro{
		cfg:            cfg,
		stopChan:       make(chan bool),
		webhookManager: webhook.NewManager(cfg),
		scanner:        scanner,
	}
}

func (m *Macro) IsRunning() bool {
	return m.running
}

func (m *Macro) Start() error {
	if m.running {
		return nil
	}

	buttons := m.cfg.MacroButtons
	if buttons["break_position"] == nil || buttons["inventory"] == nil {
		log.Println("[Macro] Not all buttons configured")
		return nil
	}

	m.running = true
	go m.run()
	return nil
}

func (m *Macro) Stop() {
	if !m.running {
		return
	}
	m.running = false
	m.stopChan <- true
}

func (m *Macro) run() {
	defer func() {
		m.running = false
		robotgo.MouseToggle("up")
	}()

	cycle := 1
	holdDuration := 5 * time.Minute
	if duration, ok := m.cfg.MacroSettings["hold_duration"].(float64); ok {
		holdDuration = time.Duration(duration) * time.Minute
	}

	autoSell := true
	if sell, ok := m.cfg.MacroSettings["auto_sell"].(bool); ok {
		autoSell = sell
	}

	for m.running {
		log.Printf("[Macro] Starting cycle %d", cycle)

		// Hold M1 at break position
		breakPos := m.cfg.MacroButtons["break_position"]
		if breakPos != nil && breakPos.X != nil && breakPos.Y != nil {
			log.Println("[Macro] Moving to break position and holding M1...")
			robotgo.Move(*breakPos.X, *breakPos.Y)
			robotgo.MouseToggle("down")

			// Hold for duration
			select {
			case <-time.After(holdDuration):
			case <-m.stopChan:
				return
			}

			robotgo.MouseToggle("up")
			time.Sleep(500 * time.Millisecond)
		}

		// Auto-sell if enabled
		if autoSell {
			if err := m.performSell(); err != nil {
				log.Printf("[Macro] Sell error: %v", err)
			}
		}

		// Send webhook update if needed
		if m.webhookManager.ShouldSendUpdate(cycle) {
			log.Println("[Macro] Sending progress update...")
			
			var stats *ocr.Stats
			if m.webhookManager.TrackStats() {
				if s, err := m.scanner.ScanForStats(nil); err == nil {
					stats = s
				}
			}
			
			if err := m.webhookManager.SendUpdate(cycle, stats); err != nil {
				log.Printf("[Macro] Webhook error: %v", err)
			}
		}

		cycle++
		time.Sleep(500 * time.Millisecond)

		select {
		case <-m.stopChan:
			return
		default:
		}
	}
}

func (m *Macro) performSell() error {
	log.Println("[Macro] Opening inventory...")

	// Open inventory (E key or click)
	invButton := m.cfg.MacroButtons["inventory"]
	if invButton != nil {
		if invButton.Key != nil {
			robotgo.KeyTap(*invButton.Key)
		} else if invButton.X != nil && invButton.Y != nil {
			robotgo.Click(*invButton.X, *invButton.Y)
		}
	}
	time.Sleep(500 * time.Millisecond)

	// Click Sell tab
	if sellTab := m.cfg.MacroButtons["sell_tab"]; sellTab != nil && sellTab.X != nil && sellTab.Y != nil {
		log.Println("[Macro] Clicking Sell tab...")
		robotgo.Click(*sellTab.X, *sellTab.Y)
		time.Sleep(300 * time.Millisecond)
	}

	// Click Select All
	if selectAll := m.cfg.MacroButtons["select_all"]; selectAll != nil && selectAll.X != nil && selectAll.Y != nil {
		log.Println("[Macro] Clicking Select All...")
		robotgo.Click(*selectAll.X, *selectAll.Y)
		time.Sleep(300 * time.Millisecond)
	}

	// Click Accept
	if accept := m.cfg.MacroButtons["accept"]; accept != nil && accept.X != nil && accept.Y != nil {
		log.Println("[Macro] Clicking Accept...")
		robotgo.Click(*accept.X, *accept.Y)
		time.Sleep(300 * time.Millisecond)
	}

	// Click Yes confirm
	if yesConfirm := m.cfg.MacroButtons["yes_confirm"]; yesConfirm != nil && yesConfirm.X != nil && yesConfirm.Y != nil {
		log.Println("[Macro] Clicking Yes...")
		robotgo.Click(*yesConfirm.X, *yesConfirm.Y)
		time.Sleep(300 * time.Millisecond)
	}

	// Close menu
	if closeMenu := m.cfg.MacroButtons["close_menu"]; closeMenu != nil && closeMenu.X != nil && closeMenu.Y != nil {
		log.Println("[Macro] Closing menu...")
		robotgo.Click(*closeMenu.X, *closeMenu.Y)
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}
