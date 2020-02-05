package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/entry", handleEntry)
	http.HandleFunc("/exit", handleExit)

	_ = http.ListenAndServe(":3000", nil)
}

type Event struct {
	Timestamp string `json:"timestamp"`
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
