package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func receiveFrames(res http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(32 << 20)
	if err != nil {
		// Error parsing form
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	fhs := req.MultipartForm.File
	for key, fh := range fhs {
		for idx, f := range fh {
			file, err := f.Open()
			if err != nil {
				// Error reading file
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			newFile, err := os.OpenFile(
				fmt.Sprintf("img/upload-%s-%d-%s", key, idx, f.Filename),
				os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)

			if err != nil {
				// Error creating file
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			defer newFile.Close()

			bytes, err := ioutil.ReadAll(file)
			if err != nil {
				// Error reading file
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}

			_, err = newFile.Write(bytes)
			if err != nil {
				// Error writing file
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	_, _ = fmt.Fprintf(res, "Saved files")
}
