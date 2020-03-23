package valet

import (
	"bytes"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"mime/multipart"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

var client *mongo.Client

func InitializeMongoClient() error {
	c, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		return err
	}

	client = c

	return nil
}

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

func Detect(timestamp time.Time, filenames []string) error {
	errChannel := make(chan error)
	args := append([]string{"src/valet/Detector.py"}, strings.Join(filenames, " "))

	go func() {
		// detect.py accepts a list of frames sorted by time
		// The number of cars in or out of the lot are returned
		// via sout, which will be redirected to storage in mongoDB
		cmd := exec.Command("python", args...)
		stdout, err := cmd.StdoutPipe()
		cmd.Stderr = os.Stderr

		if err != nil {
			errChannel <- err
			return
		}

		if err := cmd.Start(); err != nil {
			errChannel <- err
		}

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(stdout)
		if err != nil {
			errChannel <- err
			return
		}
		res := buf.String()

		collection := client.Database("Diddle_North").Collection("Log")
		_, err = collection.InsertOne(context.Background(), bson.M{
			"timestamp": timestamp,
			"entries":   res,
		})
		if err != nil {
			errChannel <- err
			return
		}
	}()

	return <-errChannel
}
