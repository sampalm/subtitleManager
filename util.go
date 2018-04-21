package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// PrintInfo print out subtitle manager logo and some infos.
func PrintInfo() {
	messageInfo := []byte("+--------------------------------------------------------------------------+\n|█▀▀ █░░█ █▀▀▄ ▀▀█▀▀ ░▀░ ▀▀█▀▀ █░░ █▀▀ █▀▀   █▀▀ ░▀░ █░░ ▀▀█▀▀ █▀▀ █▀▀█ █▀▀|\n|▀▀█ █░░█ █▀▀▄ ░░█░░ ▀█▀ ░░█░░ █░░ █▀▀ ▀▀█   █▀▀ ▀█▀ █░░ ░░█░░ █▀▀ █▄▄▀ ▀▀█|\n|▀▀▀ ░▀▀▀ ▀▀▀░ ░░▀░░ ▀▀▀ ░░▀░░ ▀▀▀ ▀▀▀ ▀▀▀   ▀░░ ▀▀▀ ▀▀▀ ░░▀░░ ▀▀▀ ▀░▀▀ ▀▀▀|\n+--------------------------------------------------------------------------+\n")

	fmt.Fprintf(os.Stdout, "%s", messageInfo)
}

// PrintHelp print out basic instructions of how to use subtitle manager.
func PrintHelp() {
	fmt.Fprintf(os.Stdout, `
Flags: 
	-p="Set the main path of the files" -e="Set the extension" -v="Set the version" -m="Set the path to move files"

Usage:
	subtitlemanager [flags]

Flags info:
	p	Execute the program with "-p" flag to set the path that contains the subtitles. eg: "./subs/" (Required)
	e	You can set the file extension. eg: ".sub" (Optional)
	v	You also can set the subtitle version. eg: "720p-WEB". (Optional)
	m	You can set the folder to where the files will be moved. eg: "./my-subs/" (Optional)
		If you set this flag no files will be deleted, only the files that matched will be moved.
	
Additional:
	If any problems occur, or you have any suggestions please send me an email: samuelpalmeira@outlook.com
	`)
}

// CreateLogFile will create or re-write a log with all error that occurs while the program has been executed.
func CreateLogFile(logs []Log) error {
	time := time.Now().Format("01/_2/2006 - 15:04:05")
	file := filepath.Join("./log.txt")
	f, err := os.OpenFile(file, syscall.O_RDWR|syscall.O_CREAT, 0777)
	if err != nil {
		return err
	}

	for _, log := range logs {
		if _, err := fmt.Fprintf(f, "**************** LOG ERROR %v ****************\r\n**** %s: %s\r\n", time, log.Func, log.Err); err != nil {
			return err
		}
	}
	return nil

}
