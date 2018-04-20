package main

import "fmt"

// PrintInfo print out subtitle manager logo and some infos.
func PrintInfo() {
	messageInfo := []byte("+--------------------------------------------------------------------------+\n|█▀▀ █░░█ █▀▀▄ ▀▀█▀▀ ░▀░ ▀▀█▀▀ █░░ █▀▀ █▀▀   █▀▀ ░▀░ █░░ ▀▀█▀▀ █▀▀ █▀▀█ █▀▀|\n|▀▀█ █░░█ █▀▀▄ ░░█░░ ▀█▀ ░░█░░ █░░ █▀▀ ▀▀█   █▀▀ ▀█▀ █░░ ░░█░░ █▀▀ █▄▄▀ ▀▀█|\n|▀▀▀ ░▀▀▀ ▀▀▀░ ░░▀░░ ▀▀▀ ░░▀░░ ▀▀▀ ▀▀▀ ▀▀▀   ▀░░ ▀▀▀ ▀▀▀ ░░▀░░ ▀▀▀ ▀░▀▀ ▀▀▀|\n+--------------------------------------------------------------------------+\n")

	fmt.Printf("%s", messageInfo)
}

// PrintHelp print out basic instructions of how to use subtitle manager.
func PrintHelp() {
	fmt.Printf(`
Flags: -p="Se the main path of the files" -e="Set the extension" -v="Set the version"
=== HOW TO USE ===
* Execute the program with "-p" flag to set the path that contains the subtitles. eg: "./subs/" (Required)
* You can set the file extension. eg: ".sub" (Optional)
* You also can set the subtitle version. eg: "720p-WEB". (Optional)
* You can set the folder to where the files will be moved. eg: "./my-subs/" (Optional)
	
If any problems occur, or you have any suggestions please send me an email: samuelpalmeira@outlook.com
	`)
}
