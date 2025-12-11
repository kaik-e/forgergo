package webui

import (
	"embed"
	"encoding/json"
	"fmt"
	"forger-companion/internal/app"
	"forger-companion/internal/config"
	"log"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

type Server struct {
	app *app.App
	cfg *config.Config
}

func NewServer(application *app.App, cfg *config.Config) *Server {
	return &Server{
		app: application,
		cfg: cfg,
	}
}

func (s *Server) Start(port int) error {
	// Serve static files
	http.Handle("/", http.FileServer(http.FS(staticFiles)))
	
	// API endpoints
	http.HandleFunc("/api/scan", s.handleScan)
	http.HandleFunc("/api/macro/toggle", s.handleMacroToggle)
	http.HandleFunc("/api/config", s.handleConfig)
	
	addr := fmt.Sprintf("localhost:%d", port)
	log.Printf("[WebUI] Starting server at http://%s", addr)
	
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement scan endpoint
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func (s *Server) handleMacroToggle(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement macro toggle
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(s.cfg)
	} else if r.Method == "POST" {
		// TODO: Update config
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "ok",
		})
	}
}
