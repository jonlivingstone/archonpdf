package mergepdf

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
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
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(inputDir)
	if err != nil {
		log.Fatalf("cannot watch directory %s: %s\n", inputDir, err)
	}

	// Event loop to listen for events
	exiting := false
	var timer *time.Timer = time.NewTimer(0)
	for !exiting {
		select {
		case event := <-watcher.Events:
			// Check if the event is file creation
			if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write {
				if filepath.Ext(event.Name) == ".pdf" {
					timer = time.NewTimer(500 * time.Millisecond)
					// os.Rename(event.Name, event.Name+".done")
					break
				}

				if filepath.Ext(event.Name) == ".done" {
					// Process the file (you can replace this with your processing logic)
					err = processFile(inputDir, mergeDir, oddPrefix, evenPrefix)
					if err != nil {
						fmt.Println(err)
						exiting = true
					}
				}
			}
		case err := <-watcher.Errors:
			log.Println("Error:", err)
		case <-timer.C:
			err = renamePdfFiles(inputDir)
			if err != nil {
				log.Fatal(err)
				return
			}
		case <-signalChan:
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
