package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Frame     string    `json:"frame"`
}

func analyzeEvent(res http.ResponseWriter, req *http.Request) {
	file, _, err := req.FormFile("media")
	defer file.Close()

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	tmp, err := ioutil.TempFile(".", "upload-*")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	defer tmp.Close()
	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = tmp.Write(fileBytes)

	msg := fmt.Sprintf("You sent an image")
	_ = json.NewEncoder(res).Encode(RouteMessage{msg})
}

type RouteTest struct {
	Message string `json:"message"`
}

type RouteMessage struct {
	Message string `json:"message"`
}

func testRoute(res http.ResponseWriter, req *http.Request) {
	msg := new(RouteTest)

	err := json.NewDecoder(req.Body).Decode(msg)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(res).Encode(RouteMessage{"Hi back!"})
}
