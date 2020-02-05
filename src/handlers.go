package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Event struct {
	Timestamp time.Time `json:"timestamp"`
	PiId      uint32    `json:"pi_id"`
	Location  string    `json:"location"`
}

func handleEntry(res http.ResponseWriter, req *http.Request) {
	entry := new(Event)

	err := json.NewDecoder(req.Body).Decode(entry)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintf(res, "Entry: %s", entry.Timestamp)
}

func handleExit(res http.ResponseWriter, req *http.Request) {
	exit := new(Event)

	err := json.NewDecoder(req.Body).Decode(exit)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintf(res, "Exit: %s", exit.Timestamp)
}
