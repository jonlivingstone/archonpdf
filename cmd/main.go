package main

import (
	archonpdf "archonpdf/internal"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {

	inputDir := flag.String("input", "/var/local/archonpdf/input", "directory containing the input pages")
	mergedDir := flag.String("merged", "/var/local/archonpdf/merged", "directory where the merged files will be written")
	flag.Parse()

	if inputDir == nil || *inputDir == "" {
		str := os.Getenv("ARCHONPDF_INPUT_DIR")
		if str == "" {
			log.Fatalf("missing parameter 'inputdir' or environment variable 'ARCHONPDF_INPUT_DIR'")
		}
		inputDir = &str
	}

	if mergedDir == nil || *mergedDir == "" {
		str := os.Getenv("ARCHONPDF_MERGED_DIR")
		if str == "" {
			log.Fatalf("missing parameter 'mergeddir' or environment variable 'ARCHONPDF_MERGED_DIR'")
		}
		mergedDir = &str
	}

	checkDirExists(*inputDir)
	checkDirExists(*mergedDir)

	cleanDirectory(*inputDir)

	fmt.Println("Starting Archon PDF daemon...")

	// Start watching the folder in a goroutine
	archonpdf.WatchFolder(*inputDir, *mergedDir, "odd", "even")

	fmt.Println("Daemon has exited.")
}

func checkDirExists(dir string) {
	fileInfo, err := os.Stat(dir)
	if os.IsNotExist(err) {
		log.Fatalf("directory %s does not exist", dir)
	}

	if !fileInfo.IsDir() {
		log.Fatalf("path %s is not a directory", dir)
	}
}

func cleanDirectory(dir string) error {

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// Allow up to 1 file in the directory
	if len(files) < 2 {
		return nil
	}

	// If there was more than 1 file let's erase eveything
	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())
		log.Printf("removing file %s", filePath)
		err = os.RemoveAll(filePath)
		if err != nil {
			log.Fatalf("failed to remove %s", filePath)
			return err
		}
	}

	return nil
}
