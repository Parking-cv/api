package main

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"parking-cv/server/src/valet"
	"time"
)

// Global count of temporary folders
var FOLDERNUM int = 0

func receiveFrames(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		// Error parsing form
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	frames := make(map[time.Time]*multipart.FileHeader)

	fhs := r.MultipartForm.File
	for timestamp, fh := range fhs {
		if len(fh) != 0 {
			http.Error(w, "Only one file should be attached to each timestamp.", http.StatusBadRequest)
			return
		}

		ts, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			// Error parsing timestamp
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		frames[ts] = fh[0]
	}

	dir := fmt.Sprintf("tmp-%d", FOLDERNUM)
	FOLDERNUM += 1

	_ = os.Mkdir(dir, 0666)

	filenames, err := valet.ReadFiles(dir, frames)
	if err != nil {
		// Error reading files
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process image in separate thread
	err = valet.Detect(filenames)
	if err != nil {
		// Error during detection
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		// Delete the temporary folder
		_ = os.RemoveAll(dir)
		FOLDERNUM -= 1
	}

	_, _ = fmt.Fprintf(w, "Processed files")
}

func testRoute(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello there!")
}
