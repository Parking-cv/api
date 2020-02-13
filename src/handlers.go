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
	var firstTimeStamp time.Time
	tsSet := false

	fhs := r.MultipartForm.File
	for timestamp, fh := range fhs {
		if len(fh) != 1 {
			http.Error(w, "Only one file should be attached to each timestamp.", http.StatusBadRequest)
			return
		}

		ts, err := time.Parse(time.RFC3339, timestamp)
		if tsSet == false {
			firstTimeStamp = ts
			tsSet = true
		}
		if err != nil {
			// Error parsing timestamp
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		frames[ts] = fh[0]
	}

	dir := fmt.Sprintf("tmp-%d", FOLDERNUM)
	FOLDERNUM += 1

	err = os.Mkdir(dir, 0777)
	if err != nil {
		// Error creating directory
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = os.Chmod(dir, 0777)
	if err != nil {
		// Error setting directory permissions
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer cleanUp(dir)

	filenames, err := valet.ReadFiles(dir, frames)
	if err != nil {
		// Error reading files
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process image in separate thread
	err = valet.Detect(firstTimeStamp, filenames)
	if err != nil {
		// Error during detection
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprintf(w, "Processed files")
}

// Remove all temporary directories and files
func cleanUp(dir string) error {
	err := os.RemoveAll(dir)
	FOLDERNUM -= 1
	return err
}

func testRoute(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello there!")
}
