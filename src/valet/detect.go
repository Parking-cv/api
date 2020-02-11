package valet

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// For sorting by time
type timeSlice []time.Time

func (t timeSlice) Len() int {
	return len(t)
}

func (t timeSlice) Less(i, j int) bool {
	return t[i].Before(t[j])
}

func (t timeSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Sort frames by timestamp and return list of filenames
func ReadFiles(dir string, frames map[time.Time]*multipart.FileHeader) ([]string, error) {
	filenames := make(map[time.Time]string)

	for timestamp, fh := range frames {
		file, err := fh.Open()
		if err != nil {
			// Error opening file
			return nil, err
		}
		defer file.Close()

		filename := fmt.Sprintf("%s/%s-%s", dir, timestamp.Format(time.RFC3339), fh.Filename)
		filenames[timestamp] = filename

		newFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			return nil, err
		}
		defer newFile.Close()

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			// Error reading file
			return nil, err
		}

		_, err = newFile.Write(bytes)
		if err != nil {
			// Error writing file
			return nil, err
		}
	}

	// Sort filenames by date and append to return value
	sortedTimestamps := make(timeSlice, 0, len(filenames))
	sortedFilenames := make([]string, 0, len(filenames))
	for t := range filenames {
		sortedTimestamps = append(sortedTimestamps, t)
	}
	sort.Sort(sortedTimestamps)

	for _, t := range sortedTimestamps {
		sortedFilenames = append(sortedFilenames, filenames[t])
	}

	return sortedFilenames, nil
}

func Detect(filenames []string) error {

	errChannel := make(chan error)
	args := append([]string{"src/valet/Detector.py"}, strings.Join(filenames, " "))

	go func() {
		// TODO set these environment variables somewhere else
		os.Setenv("PYTHONPATH", "/Users/davidwiles/PycharmProjects/parking-cv-python/venv/lib/python3.7/site-packages/:/Users/davidwiles/PycharmProjects/parking-cv-python")
		os.Setenv("PYTHON_BIN","/Users/davidwiles/PycharmProjects/parking-cv-python/venv/bin/python")
		pythonBin := os.Getenv("PYTHON_BIN")
		// detect.py accepts a list of frames sorted by time
		cmd := exec.Command(pythonBin, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		errChannel <- cmd.Run()
	}()

	return <-errChannel
}
