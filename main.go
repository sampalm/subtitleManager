package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var subname, subpath string
	var subFiles = []string{}

	// SHOW info about SubtitlesFilter
	fmt.Println("+--------------------------------------------------------------------------+")
	fmt.Println("|█▀▀ █░░█ █▀▀▄ ▀▀█▀▀ ░▀░ ▀▀█▀▀ █░░ █▀▀ █▀▀   █▀▀ ░▀░ █░░ ▀▀█▀▀ █▀▀ █▀▀█ █▀▀|")
	fmt.Println("|▀▀█ █░░█ █▀▀▄ ░░█░░ ▀█▀ ░░█░░ █░░ █▀▀ ▀▀█   █▀▀ ▀█▀ █░░ ░░█░░ █▀▀ █▄▄▀ ▀▀█|")
	fmt.Println("|▀▀▀ ░▀▀▀ ▀▀▀░ ░░▀░░ ▀▀▀ ░░▀░░ ▀▀▀ ▀▀▀ ▀▀▀   ▀░░ ▀▀▀ ▀▀▀ ░░▀░░ ▀▀▀ ▀░▀▀ ▀▀▀|")
	fmt.Println("+--------------------------------------------------------------------------+")
	fmt.Println(`ABOUT: SubtitlesFilter can make your life more easier by selecting only the subtitles
	you need to watch your movies, series or whatever you want.
	************** HOW TO USE **************
	Put the address where your subtitles are and after that just enter the version of the subtitle
	you want, after the program finish the task your subtitles will be in the backup folder that 
	will be created automatically. Enjoy ;) 
	`)
	

	// GET Subtitles path
	fmt.Print("Enter the subtitle path[Ex: ./subtitles/]: ")
	fmt.Scan(&subpath)

	// GET Subtitle version
	fmt.Print("Enter the subtitle version[Ex: HDTV.x264]: ")
	fmt.Scan(&subname)

	// SEARCH files in the address
	if subpath != "" && subname != "" {
		subFiles = search(subpath, subname)
	} else {
		fmt.Println("Invalid subtitle path or name.")
		time.Sleep(time.Duration(3 * time.Second))
	}

	// FILTER the files in the address
	// DELETE the files that dont match
	if filter(subpath, subFiles); deleteFiles(subpath) {
		// Finish program
		fmt.Println("+-------------------------------------------------+")
		fmt.Println("|**** TASK COMPLETED ")
		fmt.Println("|**** YOUR SUBTITLES ARE IN THE BACKUP FOLDER")
		fmt.Println("|**** " + subpath + "backup/ ")
		fmt.Println("+-------------------------------------------------+")
		time.Sleep(time.Duration(5 * time.Second))
	}
}

func search(path, sub string) []string {
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
		return subFiles
	} else {
		log.Println("None files match with these description.")
		time.Sleep(time.Duration(5 * time.Second))
		os.Exit(1)
	}
	return nil
}

func filter(path string, files []string) bool {
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

	return true

}

func deleteFiles(path string) bool {
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

	return true

}
