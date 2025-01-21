package mergepdf

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/jhenstridge/go-inotify"
)

func processFile(inputDir, mergeDir, oddPrefix, evenPrefix string) error {
	files, err := os.ReadDir(inputDir)
	if err != nil {
		return err
	}

	oddFileName := findFileWithPrefix(files, oddPrefix)
	evenFileName := findFileWithPrefix(files, evenPrefix)

	if oddFileName == "" || evenFileName == "" {
		return nil
	}

	oddFilePath := filepath.Join(inputDir, oddFileName)
	evenFilePath := filepath.Join(inputDir, evenFileName)

	resultFileName := strings.Replace(time.Now().Format(time.RFC3339)+".pdf", ":", "-", -1)
	resultFilePath := filepath.Join(mergeDir, resultFileName)

	fmt.Printf("merging odd pages file %s with even pages file %s into %s\n",
		oddFileName, evenFileName, resultFilePath)

	err = mergeOddEvenPdfs(oddFilePath, evenFilePath, resultFilePath)
	if err != nil {
		return err
	}

	fmt.Printf("files merged into %s\n", resultFilePath)

	err = os.Remove(oddFilePath)
	if err != nil {
		return err
	}

	err = os.Remove(evenFilePath)
	if err != nil {
		return err
	}

	return nil
}

func findFileWithPrefix(files []os.DirEntry, prefix string) string {
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasPrefix(file.Name(), prefix) && strings.HasSuffix(file.Name(), ".pdf.done") {
			return file.Name()
		}
	}

	return ""
}

func WatchFolder(inputDir, mergeDir, oddPrefix, evenPrefix string) {

	// Create a channel to handle SIGTERM signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a new watcher
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	_, err = watcher.AddWatch(inputDir, inotify.IN_ALL_EVENTS)
	if err != nil {
		log.Fatalf("cannot watch directory %s: %s\n", inputDir, err)
	}

	// Event loop to listen for events
	exiting := false
	const infinite = math.MaxInt32 * time.Second
	var timer *time.Timer = time.NewTimer(infinite)
	for !exiting {
		select {
		case event := <-watcher.Event:
			if event.Mask&inotify.IN_CLOSE_WRITE == inotify.IN_CLOSE_WRITE {
				if filepath.Ext(event.Name) == ".pdf" {
					timer = time.NewTimer(3 * time.Second)
				}
			}
		case <-timer.C:
			renamePdfFiles(inputDir)
			err = processFile(inputDir, mergeDir, oddPrefix, evenPrefix)
			if err != nil {
				fmt.Println(err)
			}

		case <-signalChan:
			watcher.Close()
			exiting = true
		}
	}
}

func renamePdfFiles(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".pdf") {
			src := filepath.Join(dir, file.Name())
			dest := src + ".done"
			err = os.Rename(src, dest)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
