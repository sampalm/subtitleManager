package main

import (
	"time"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var subname, subpath string

	// GET Subtitles path
	fmt.Print("Enter the subtitle path[Ex: ./subtitles/]: ")
	fmt.Scan(&subpath)

	// GET Subtitle version
	fmt.Print("Enter the subtitle version[Ex: HDTV.x264]: ")
	fmt.Scan(&subname)

	if subpath != "" && subname != "" {
		search(subpath, subname)
	} else {
		fmt.Println("Invalid subtitle path or name.")
		time.Sleep(time.Duration(3 * time.Second))
	}
}

func search(path, sub string) {
	var subFiles = []string{}

	files, err := filepath.Glob(path + "*" + sub + "*.srt")
	if err != nil {
		log.Println(err)
		time.Sleep(time.Duration(5 * time.Second))
		os.Exit(1)
	}

	for _, f := range files {
		subFiles = append(subFiles, f)
	}

	if subFiles != nil && len(subFiles) > 0 {
		filter(path, subFiles)
	} else {
		log.Println("None files match with these description.")
		time.Sleep(time.Duration(5 * time.Second))
		os.Exit(1)
	}
}

func filter(path string, files []string) {
	var subs string
	var bkp = path + "backup/"

	// Create backup path
	if _, err := os.Stat(bkp); os.IsNotExist(err) {
		os.Mkdir(bkp, 0644)
	}

	// Move files to backup
	for _, f := range files {
		subs = strings.Split(f, "\\")[1]
		err := os.Rename(f, path+"backup/"+subs)

		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(5 * time.Second))
			os.Exit(2)
		}
	}

	// Delete all other files
	deleteFiles(path)
}

func deleteFiles(path string) {
	files, err := filepath.Glob(path + "*.srt")
	if err != nil {
		log.Println(err)
		time.Sleep(time.Duration(5 * time.Second))
		os.Exit(3)
	}

	for _, f := range files {
		err := os.Remove(f)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(5 * time.Second))
			os.Exit(4)
		}
	}

	// Finish program
	fmt.Println("+-------------------------------------------------+")
	fmt.Println("|**************** TASK COMPLETED *****************|")
	fmt.Println("|**** YOUR SUBTITLES ARE IN THE BACKUP FOLDER ****|")
	fmt.Println("|********* "+path+"backup/ ********|")
	fmt.Println("+-------------------------------------------------+")
	time.Sleep(time.Duration(5 * time.Second))
}
