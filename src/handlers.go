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
	Image     string    `json:"filename"`
}

func analyzeEvent(res http.ResponseWriter, req *http.Request) {
	entry := new(Event)

	err := json.NewDecoder(req.Body).Decode(entry)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	//err = DB.Ping()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprintf(res, "Entry: %s", entry.Timestamp)
}

type TestEvent struct {
	Message string `json:"message"`
}

func testRoute(res http.ResponseWriter, req *http.Request) {
	msg := new(TestEvent)

	err := json.NewDecoder(req.Body).Decode(msg)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprintf(res, "You sent: %s", msg.Message)
}
