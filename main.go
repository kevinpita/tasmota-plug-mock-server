// Package main implements a simple mock server for a Tasmota smart plug.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	energyTotalToday     = 2.200
	energyTotalYesterday = 2.000
	energyTotalTotal     = 5000.920
)

// Server represents the mock server.
type Server struct {
	mu         sync.Mutex
	powerState string
}

// NewServer creates a new mock server.
func NewServer() *Server {
	return &Server{
		powerState: "OFF",
		mu:         sync.Mutex{},
	}
}

func main() {
	const readHeaderTimeout = 3 * time.Second
	server := NewServer()
	http.HandleFunc("/cm", server.handleCmnd)
	log.Println("Starting server on :8080")
	srv := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: readHeaderTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("could not listen on port 8080 %v", err)
	}
}

func (s *Server) handleCmnd(w http.ResponseWriter, r *http.Request) {
	cmnd := r.URL.Query().Get("cmnd")
	switch cmnd {
	case "EnergyTotal":
		s.handleEnergyTotal(w)
	case "Power On":
		s.handlePowerOn(w)
	case "Power Off":
		s.handlePowerOff(w)
	case "Power":
		s.handlePowerStatus(w)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) handleEnergyTotal(w http.ResponseWriter) {
	response := map[string]interface{}{
		"EnergyTotal": map[string]interface{}{
			"Total":     energyTotalTotal,
			"Yesterday": energyTotalYesterday,
			"Today":     energyTotalToday,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "could not encode json", http.StatusInternalServerError)
	}
}

func (s *Server) handlePowerOn(w http.ResponseWriter) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.powerState = "ON"

	response := map[string]string{"POWER": "ON"}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "could not encode json", http.StatusInternalServerError)
	}
}

func (s *Server) handlePowerOff(w http.ResponseWriter) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.powerState = "OFF"

	response := map[string]string{"POWER": "OFF"}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "could not encode json", http.StatusInternalServerError)
	}
}

func (s *Server) handlePowerStatus(w http.ResponseWriter) {
	s.mu.Lock()
	defer s.mu.Unlock()

	response := map[string]string{"POWER": s.powerState}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "could not encode json", http.StatusInternalServerError)
	}
}
