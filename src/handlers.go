package main

import (
	"encoding/json"
	"net/http"
)

// Payload - структура входящего JSON
type Payload struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// handleInput - POST /data
func handleInput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var p Payload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	// Отправляем в канал (неблокирующая попытка)
	select {
	case dataChan <- p.Value:
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status":"queued"}`))
	default:
		// Если воркер не справляется и очередь забита
		http.Error(w, "Queue full", http.StatusServiceUnavailable)
	}
}

// handleHealth - GET /health
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}